package router

import (
	"context"
	"sync"
	"time"
)

const (
	IsolateModeShared    = "shared"
	IsolateModeSchema    = "schema"
	IsolateModeDedicated = "dedicated"
)

type TenantRouteInfo struct {
	TenantID    string
	IsolateMode string
	DSN         string
	MaxConns    int32
}

type TenantResolver interface {
	ResolveTenant(ctx context.Context, tenantID string) (*TenantRouteInfo, error)
}

type cachedItem struct {
	info      *TenantRouteInfo
	expiresAt time.Time
}

type CachedTenantResolver struct {
	underlying TenantResolver
	ttl        time.Duration
	mu         sync.RWMutex
	cache      map[string]cachedItem
}

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
