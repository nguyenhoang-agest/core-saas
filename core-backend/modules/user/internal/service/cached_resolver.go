package service

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/eventbus"
	"github.com/example/core-backend/pkg/rbac"
	"gorm.io/gorm"
)

const (
	AuthzInvalidatedTopic = "iam.authz.invalidated"
	DefaultCacheTTL       = 15 * time.Second
)

type cacheEntry struct {
	permissions  []string
	authzVersion int64
	expiresAt    time.Time
}

// CachedPermissionResolver implements rbac.PermissionResolver with tiered in-memory caching,
// 15s TTL, authz_version verification, and eventbus-driven invalidation.
type CachedPermissionResolver struct {
	db  *gorm.DB
	bus *eventbus.Bus

	mu    sync.RWMutex
	cache map[string]*cacheEntry // key: tenantID + ":" + userID
}

// NewCachedPermissionResolver constructs a CachedPermissionResolver and binds event bus listeners.
func NewCachedPermissionResolver(db *gorm.DB, bus *eventbus.Bus) *CachedPermissionResolver {
	resolver := &CachedPermissionResolver{
		db:    db,
		bus:   bus,
		cache: make(map[string]*cacheEntry),
	}

	if bus != nil {
		eventbus.Subscribe(bus, AuthzInvalidatedTopic, func(ctx context.Context, payload map[string]string) error {
			tenantID := payload["tenant_id"]
			userID := payload["user_id"]
			resolver.Invalidate(ctx, tenantID, userID)
			return nil
		})
	}

	return resolver
}

func cacheKey(tenantID, userID string) string {
	return tenantID + ":" + userID
}

// ResolvePermissions computes the complete set of concrete permissions granted to a user.
func (r *CachedPermissionResolver) ResolvePermissions(ctx context.Context, tenantID, userID string, isHost bool) ([]string, error) {
	key := cacheKey(tenantID, userID)
	now := time.Now()

	// 1. Fast path: check in-memory cache
	r.mu.RLock()
	entry, found := r.cache[key]
	r.mu.RUnlock()

	if found && now.Before(entry.expiresAt) {
		return entry.permissions, nil
	}

	// 2. Fetch current authz_version for the tenant
	var currentVersion int64 = 1
	var authzVer model.AuthzVersion
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).First(&authzVer).Error; err == nil {
		currentVersion = authzVer.Version
	}

	// If entry exists and its authzVersion matches current DB version, refresh TTL and return
	if found && entry.authzVersion == currentVersion {
		r.mu.Lock()
		entry.expiresAt = now.Add(DefaultCacheTTL)
		r.mu.Unlock()
		return entry.permissions, nil
	}

	// 3. Cache miss or version changed: query database
	var permKeys []string
	query := `
		SELECT DISTINCT rp.permission_key
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		WHERE ur.tenant_id = ? AND ur.user_id = ?
	`
	if err := r.db.WithContext(ctx).Raw(query, tenantID, userID).Scan(&permKeys).Error; err != nil {
		return nil, fmt.Errorf("resolving permissions from DB for user %s: %w", userID, err)
	}

	// If no user_roles exist, check legacy users.role as fallback during transition
	if len(permKeys) == 0 {
		var user model.User
		if err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, userID).First(&user).Error; err == nil && user.Role != "" {
			permKeys = rbac.GetRegistry().ExpandPermissions(rbac.GetDefaultRolePermissions(rbac.Role(user.Role)), isHost)
		}
	}

	// Filter out ScopePlatform permissions if not host tenant
	reg := rbac.GetRegistry()
	filtered := make([]string, 0, len(permKeys))
	for _, k := range permKeys {
		if perm, ok := reg.Get(k); ok {
			if perm.Scope == rbac.ScopeTenant || isHost {
				filtered = append(filtered, k)
			}
		}
	}
	sort.Strings(filtered)

	// 4. Update in-memory cache
	r.mu.Lock()
	r.cache[key] = &cacheEntry{
		permissions:  filtered,
		authzVersion: currentVersion,
		expiresAt:    now.Add(DefaultCacheTTL),
	}
	r.mu.Unlock()

	return filtered, nil
}

// HasPermission checks if the user possesses the specified concrete permission.
func (r *CachedPermissionResolver) HasPermission(ctx context.Context, tenantID, userID string, isHost bool, requiredKey string) (bool, error) {
	perms, err := r.ResolvePermissions(ctx, tenantID, userID, isHost)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p == requiredKey || p == "*" {
			return true, nil
		}
	}
	return false, nil
}

// Invalidate removes cached entries for a user or entire tenant.
func (r *CachedPermissionResolver) Invalidate(ctx context.Context, tenantID, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if userID != "" {
		delete(r.cache, cacheKey(tenantID, userID))
		return
	}

	// Invalidate all users in this tenant
	prefix := tenantID + ":"
	for k := range r.cache {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(r.cache, k)
		}
	}
}
