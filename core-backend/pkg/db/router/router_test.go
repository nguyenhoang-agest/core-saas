package router_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/core-backend/pkg/db/router"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mockResolver struct {
	resolveFunc func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error)
	callCount   int32
}

func (m *mockResolver) ResolveTenant(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
	atomic.AddInt32(&m.callCount, 1)
	if m.resolveFunc != nil {
		return m.resolveFunc(ctx, tenantID)
	}
	return nil, nil
}

func TestStaticDBProvider(t *testing.T) {
	dummyPool := &pgxpool.Pool{}
	provider := router.NewStaticDBProvider(dummyPool)

	pool, err := provider.GetPool(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool != dummyPool {
		t.Fatalf("expected pool to match dummyPool")
	}

	if provider.GetDefaultPool() != dummyPool {
		t.Fatalf("expected GetDefaultPool to match dummyPool")
	}
}

func TestDynamicDBRouter_NoTenantInContext(t *testing.T) {
	defaultPool := &pgxpool.Pool{}
	res := &mockResolver{}
	r := router.NewDynamicDBRouter(
		defaultPool,
		res,
		router.WithIdleTimeout(0),
		router.WithCloseFunc(func(p *pgxpool.Pool) {}),
	)
	defer r.Close()

	pool, err := r.GetPool(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool != defaultPool {
		t.Fatalf("expected defaultPool, got %v", pool)
	}
	if atomic.LoadInt32(&res.callCount) != 0 {
		t.Fatalf("resolver should not have been called")
	}
}

func TestDynamicDBRouter_SharedTenant(t *testing.T) {
	defaultPool := &pgxpool.Pool{}
	res := &mockResolver{
		resolveFunc: func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
			return &router.TenantRouteInfo{
				TenantID:    tenantID,
				IsolateMode: router.IsolateModeShared,
			}, nil
		},
	}
	r := router.NewDynamicDBRouter(
		defaultPool,
		res,
		router.WithIdleTimeout(0),
		router.WithCloseFunc(func(p *pgxpool.Pool) {}),
	)
	defer r.Close()

	ctx := tenantctx.WithContext(context.Background(), "tenant-shared-1")
	pool, err := r.GetPool(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool != defaultPool {
		t.Fatalf("expected defaultPool for shared tenant, got %v", pool)
	}
}

func TestDynamicDBRouter_DedicatedTenant(t *testing.T) {
	defaultPool := &pgxpool.Pool{}
	dedicatedPool := &pgxpool.Pool{}

	var connectCalls int32
	mockConnect := func(ctx context.Context, dsn string, maxConns int32) (*pgxpool.Pool, error) {
		atomic.AddInt32(&connectCalls, 1)
		if dsn != "postgres://test_dedicated_dsn" {
			t.Errorf("unexpected dsn: %s", dsn)
		}
		if maxConns != 15 {
			t.Errorf("unexpected maxConns: %d", maxConns)
		}
		return dedicatedPool, nil
	}

	var closeCalls int32
	mockClose := func(p *pgxpool.Pool) {
		atomic.AddInt32(&closeCalls, 1)
	}

	res := &mockResolver{
		resolveFunc: func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
			return &router.TenantRouteInfo{
				TenantID:    tenantID,
				IsolateMode: router.IsolateModeDedicated,
				DSN:         "postgres://test_dedicated_dsn",
				MaxConns:    15,
			}, nil
		},
	}

	r := router.NewDynamicDBRouter(
		defaultPool,
		res,
		router.WithConnectFunc(mockConnect),
		router.WithCloseFunc(mockClose),
		router.WithIdleTimeout(0),
	)

	ctx := tenantctx.WithContext(context.Background(), "tenant-dedicated-1")

	// First call should connect
	pool1, err := r.GetPool(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool1 != dedicatedPool {
		t.Fatalf("expected dedicatedPool, got %v", pool1)
	}
	if atomic.LoadInt32(&connectCalls) != 1 {
		t.Fatalf("expected 1 connect call, got %d", connectCalls)
	}

	// Second call should return cached pool
	pool2, err := r.GetPool(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool2 != dedicatedPool {
		t.Fatalf("expected cached dedicatedPool, got %v", pool2)
	}
	if atomic.LoadInt32(&connectCalls) != 1 {
		t.Fatalf("expected connect call to remain 1, got %d", connectCalls)
	}

	r.Close()
	if atomic.LoadInt32(&closeCalls) != 1 {
		t.Fatalf("expected 1 close call when router closed, got %d", closeCalls)
	}
}

func TestDynamicDBRouter_DedicatedMissingDSN(t *testing.T) {
	defaultPool := &pgxpool.Pool{}
	res := &mockResolver{
		resolveFunc: func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
			return &router.TenantRouteInfo{
				TenantID:    tenantID,
				IsolateMode: router.IsolateModeDedicated,
				DSN:         "",
			}, nil
		},
	}
	r := router.NewDynamicDBRouter(
		defaultPool,
		res,
		router.WithIdleTimeout(0),
		router.WithCloseFunc(func(p *pgxpool.Pool) {}),
	)
	defer r.Close()

	ctx := tenantctx.WithContext(context.Background(), "tenant-bad-dsn")
	_, err := r.GetPool(ctx)
	if err == nil {
		t.Fatalf("expected error for dedicated tenant with empty dsn")
	}
}

func TestCachedTenantResolver(t *testing.T) {
	underlying := &mockResolver{
		resolveFunc: func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
			return &router.TenantRouteInfo{
				TenantID:    tenantID,
				IsolateMode: router.IsolateModeShared,
			}, nil
		},
	}

	cached := router.NewCachedTenantResolver(underlying, 100*time.Millisecond)

	ctx := context.Background()

	// Call 1
	info1, err := cached.ResolveTenant(ctx, "t-1")
	if err != nil || info1.TenantID != "t-1" {
		t.Fatalf("failed call 1: %v", err)
	}
	if atomic.LoadInt32(&underlying.callCount) != 1 {
		t.Fatalf("expected 1 underlying call, got %d", underlying.callCount)
	}

	// Call 2 within TTL
	info2, err := cached.ResolveTenant(ctx, "t-1")
	if err != nil || info2.TenantID != "t-1" {
		t.Fatalf("failed call 2: %v", err)
	}
	if atomic.LoadInt32(&underlying.callCount) != 1 {
		t.Fatalf("expected cached result, but underlying was called %d times", underlying.callCount)
	}

	// Invalidate manually
	cached.Invalidate("t-1")

	info3, err := cached.ResolveTenant(ctx, "t-1")
	if err != nil || info3.TenantID != "t-1" {
		t.Fatalf("failed call 3: %v", err)
	}
	if atomic.LoadInt32(&underlying.callCount) != 2 {
		t.Fatalf("expected 2 underlying calls after invalidate, got %d", underlying.callCount)
	}
}

func TestDynamicDBRouter_ResolverError(t *testing.T) {
	defaultPool := &pgxpool.Pool{}
	res := &mockResolver{
		resolveFunc: func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
			return nil, errors.New("database lookup error")
		},
	}
	r := router.NewDynamicDBRouter(
		defaultPool,
		res,
		router.WithIdleTimeout(0),
		router.WithCloseFunc(func(p *pgxpool.Pool) {}),
	)
	defer r.Close()

	ctx := tenantctx.WithContext(context.Background(), "t-error")
	_, err := r.GetPool(ctx)
	if err == nil {
		t.Fatalf("expected error from resolver failure, got nil")
	}
}

func TestDynamicDBRouter_ReapIdlePools(t *testing.T) {
	defaultPool := &pgxpool.Pool{}
	dedicatedPool := &pgxpool.Pool{}

	var closeCalls int32
	mockClose := func(p *pgxpool.Pool) {
		atomic.AddInt32(&closeCalls, 1)
	}

	res := &mockResolver{
		resolveFunc: func(ctx context.Context, tenantID string) (*router.TenantRouteInfo, error) {
			return &router.TenantRouteInfo{
				TenantID:    tenantID,
				IsolateMode: router.IsolateModeDedicated,
				DSN:         "postgres://test_dedicated_dsn",
				MaxConns:    10,
			}, nil
		},
	}

	// 50ms idle timeout
	r := router.NewDynamicDBRouter(
		defaultPool,
		res,
		router.WithConnectFunc(func(ctx context.Context, dsn string, maxConns int32) (*pgxpool.Pool, error) {
			return dedicatedPool, nil
		}),
		router.WithCloseFunc(mockClose),
		router.WithIdleTimeout(50*time.Millisecond),
	)
	defer r.Close()

	ctx := tenantctx.WithContext(context.Background(), "t-reap")
	_, err := r.GetPool(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Wait for reaper to run (ticker runs every idleTimeout / 2 = 25ms, wait ~120ms)
	time.Sleep(150 * time.Millisecond)

	if atomic.LoadInt32(&closeCalls) < 1 {
		t.Fatalf("expected idle pool to be reaped and closed")
	}
}

