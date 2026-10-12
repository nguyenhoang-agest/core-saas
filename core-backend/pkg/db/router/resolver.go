package router

import (
	"context"
	"sync"
	"time"
)

// Standard isolation modes
const (
	IsolateModeShared    = "shared"
	IsolateModeSchema    = "schema"
	IsolateModeDedicated = "dedicated"
)

// TenantRouteInfo carries routing and connection attributes for a tenant.
type TenantRouteInfo struct {
	TenantID    string
	IsolateMode string
	DSN         string
	MaxConns    int32
}

// TenantResolver resolves routing information for a given tenant ID.
type TenantResolver interface {
	ResolveTenant(ctx context.Context, tenantID string) (*TenantRouteInfo, error)
}

type cachedItem struct {
	info      *TenantRouteInfo
	expiresAt time.Time
}

// CachedTenantResolver wraps another TenantResolver with in-memory TTL caching.
type CachedTenantResolver struct {
	underlying TenantResolver
	ttl        time.Duration
	mu         sync.RWMutex
	cache      map[string]cachedItem
}

// NewCachedTenantResolver constructs a CachedTenantResolver with the specified TTL.
func NewCachedTenantResolver(underlying TenantResolver, ttl time.Duration) *CachedTenantResolver {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &CachedTenantResolver{
		underlying: underlying,
		ttl:        ttl,
		cache:      make(map[string]cachedItem),
	}
}

// ResolveTenant returns cached routing info if valid, otherwise delegates to underlying.
func (c *CachedTenantResolver) ResolveTenant(ctx context.Context, tenantID string) (*TenantRouteInfo, error) {
	now := time.Now()

	c.mu.RLock()
	item, ok := c.cache[tenantID]
	c.mu.RUnlock()

	if ok && now.Before(item.expiresAt) {
		return item.info, nil
	}

	info, err := c.underlying.ResolveTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[tenantID] = cachedItem{
		info:      info,
		expiresAt: now.Add(c.ttl),
	}
	c.mu.Unlock()

	return info, nil
}

// Invalidate purges the cache entry for a tenant (e.g. after config update).
func (c *CachedTenantResolver) Invalidate(tenantID string) {
	c.mu.Lock()
	delete(c.cache, tenantID)
	c.mu.Unlock()
}
