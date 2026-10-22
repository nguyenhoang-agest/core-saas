package rbac

import (
	"context"
	"sync"
)

// PermissionResolver defines the contract for resolving and checking user permissions dynamically.
type PermissionResolver interface {
	// ResolvePermissions computes the complete set of concrete permissions granted to a user.
	ResolvePermissions(ctx context.Context, tenantID, userID string, isHost bool) ([]string, error)

	// HasPermission checks if the user possesses the specified concrete permission.
	HasPermission(ctx context.Context, tenantID, userID string, isHost bool, requiredKey string) (bool, error)

	// Invalidate removes cached entries for a user or tenant.
	Invalidate(ctx context.Context, tenantID, userID string)
}

var (
	resolverMu     sync.RWMutex
	activeResolver PermissionResolver
)

// SetPermissionResolver registers the runtime dynamic permission resolver.
func SetPermissionResolver(r PermissionResolver) {
	resolverMu.Lock()
	defer resolverMu.Unlock()
	activeResolver = r
}

// GetPermissionResolver retrieves the currently active runtime permission resolver.
func GetPermissionResolver() PermissionResolver {
	resolverMu.RLock()
	defer resolverMu.RUnlock()
	return activeResolver
}
