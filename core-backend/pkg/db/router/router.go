package router

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBProvider interface {
	GetPool(ctx context.Context) (*pgxpool.Pool, error)
	GetDefaultPool() *pgxpool.Pool
	Close()
}

type StaticDBProvider struct {
	pool *pgxpool.Pool
}

func NewStaticDBProvider(pool *pgxpool.Pool) *StaticDBProvider {
	return &StaticDBProvider{pool: pool}
}

func (s *StaticDBProvider) GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	return s.pool, nil
}

func (s *StaticDBProvider) GetDefaultPool() *pgxpool.Pool {
	return s.pool
}

func (s *StaticDBProvider) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// PoolConnectFunc defines the signature for establishing a new connection pool.
type PoolConnectFunc func(ctx context.Context, dsn string, maxConns int32) (*pgxpool.Pool, error)

// PoolCloseFunc defines the signature for closing a connection pool.
type PoolCloseFunc func(pool *pgxpool.Pool)

// Option configures DynamicDBRouter.
type Option func(*DynamicDBRouter)

// WithConnectFunc overrides the pool creation function.
func WithConnectFunc(fn PoolConnectFunc) Option {
	return func(r *DynamicDBRouter) {
		r.connectFunc = fn
	}
}

// WithCloseFunc overrides the pool closing function.
func WithCloseFunc(fn PoolCloseFunc) Option {
	return func(r *DynamicDBRouter) {
		r.closeFunc = fn
	}
}

// WithIdleTimeout sets the duration after which an idle dedicated pool is evicted.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(r *DynamicDBRouter) {
		r.idleTimeout = timeout
	}
}

// DynamicDBRouter manages routing database requests to appropriate tenant connection pools.
type DynamicDBRouter struct {
	defaultPool *pgxpool.Pool
	resolver    TenantResolver
	connectFunc PoolConnectFunc
	closeFunc   PoolCloseFunc
	idleTimeout time.Duration

	mu       sync.RWMutex
	pools    map[string]*pgxpool.Pool
	lastUsed map[string]time.Time

	stopReaper chan struct{}
	reaperWG   sync.WaitGroup
}

// NewDynamicDBRouter initializes a dynamic database router.
func NewDynamicDBRouter(defaultPool *pgxpool.Pool, resolver TenantResolver, opts ...Option) *DynamicDBRouter {
	r := &DynamicDBRouter{
		defaultPool: defaultPool,
		resolver:    resolver,
		connectFunc: db.Connect,
		closeFunc: func(p *pgxpool.Pool) {
			if p != nil {
				p.Close()
			}
		},
		idleTimeout: 30 * time.Minute,
		pools:       make(map[string]*pgxpool.Pool),
		lastUsed:    make(map[string]time.Time),
		stopReaper:  make(chan struct{}),
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.idleTimeout > 0 {
		r.reaperWG.Add(1)
		go r.runIdleReaper()
	}

	return r
}

// GetPool resolves and returns the pool for the active tenant context.
func (r *DynamicDBRouter) GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	tenantID, ok := tenantctx.FromContext(ctx)
	if !ok || tenantID == "" || r.resolver == nil {
		return r.defaultPool, nil
	}

	routeInfo, err := r.resolver.ResolveTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("resolving tenant %q routing: %w", tenantID, err)
	}

	if routeInfo == nil || routeInfo.IsolateMode != IsolateModeDedicated {
		return r.defaultPool, nil
	}

	if routeInfo.DSN == "" {
		return nil, fmt.Errorf("dedicated tenant %q missing database connection dsn", tenantID)
	}

	// Fast path with read lock
	r.mu.RLock()
	pool, exists := r.pools[tenantID]
	r.mu.RUnlock()

	if exists {
		r.mu.Lock()
		r.lastUsed[tenantID] = time.Now()
		r.mu.Unlock()
		return pool, nil
	}

	// Slow path with write lock
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double check after acquiring write lock
	if pool, exists = r.pools[tenantID]; exists {
		r.lastUsed[tenantID] = time.Now()
		return pool, nil
	}

	maxConns := routeInfo.MaxConns
	if maxConns <= 0 {
		maxConns = 10
	}

	newPool, err := r.connectFunc(ctx, routeInfo.DSN, maxConns)
	if err != nil {
		return nil, fmt.Errorf("connecting dedicated pool for tenant %q: %w", tenantID, err)
	}

	r.pools[tenantID] = newPool
	r.lastUsed[tenantID] = time.Now()

	return newPool, nil
}

// GetDefaultPool returns the default/control-plane pool.
func (r *DynamicDBRouter) GetDefaultPool() *pgxpool.Pool {
	return r.defaultPool
}

// runIdleReaper periodically reaps idle tenant connection pools.
func (r *DynamicDBRouter) runIdleReaper() {
	defer r.reaperWG.Done()

	reapInterval := r.idleTimeout / 2
	if reapInterval < 10*time.Millisecond {
		reapInterval = 10 * time.Millisecond
	}

	ticker := time.NewTicker(reapInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.stopReaper:
			return
		case now := <-ticker.C:
			r.reapIdlePools(now)
		}
	}
}

func (r *DynamicDBRouter) reapIdlePools(now time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for tenantID, pool := range r.pools {
		lastTime, ok := r.lastUsed[tenantID]
		if ok && now.Sub(lastTime) > r.idleTimeout {
			r.closeFunc(pool)
			delete(r.pools, tenantID)
			delete(r.lastUsed, tenantID)
		}
	}
}

// Close gracefully closes all dynamically opened pools and stops background reaper.
func (r *DynamicDBRouter) Close() {
	select {
	case <-r.stopReaper:
		// already closed
	default:
		close(r.stopReaper)
		r.reaperWG.Wait()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for tenantID, pool := range r.pools {
		r.closeFunc(pool)
		delete(r.pools, tenantID)
	}
}
