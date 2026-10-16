package tenant_test

import (
	"context"
	"testing"

	"github.com/example/core-backend/modules/tenant"
	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/pkg/db/router"
	"github.com/example/core-backend/pkg/httpx"
)

type mockTenantService struct {
	getByIDFunc func(ctx context.Context, id string) (*model.Tenant, error)
}

func (m *mockTenantService) Create(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error) {
	return nil, nil
}

func (m *mockTenantService) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTenantService) GetAllTenant(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.Tenant], error) {
	return nil, nil
}

func TestTenantResolver(t *testing.T) {
	svc := &mockTenantService{
		getByIDFunc: func(ctx context.Context, id string) (*model.Tenant, error) {
			if id == "dedicated-t" {
				return &model.Tenant{
					BaseModel:   model.BaseModel{Id: "dedicated-t"},
					Name:        "Dedicated",
					IsolateMode: model.IsolateModeDedicated,
					DBConfig: &model.TenantDBConfig{
						DSN:      "postgres://custom_host/db",
						MaxConns: 25,
					},
				}, nil
			}
			return &model.Tenant{
				BaseModel:   model.BaseModel{Id: "shared-t"},
				Name:        "Shared",
				IsolateMode: model.IsolateModeShared,
			}, nil
		},
	}

	resolver := tenant.NewTenantResolver(svc)

	// Test dedicated resolution
	info1, err := resolver.ResolveTenant(context.Background(), "dedicated-t")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info1.IsolateMode != router.IsolateModeDedicated {
		t.Fatalf("expected dedicated mode, got %s", info1.IsolateMode)
	}
	if info1.DSN != "postgres://custom_host/db" {
		t.Fatalf("unexpected dsn: %s", info1.DSN)
	}
	if info1.MaxConns != 25 {
		t.Fatalf("unexpected maxConns: %d", info1.MaxConns)
	}

	// Test shared resolution
	info2, err := resolver.ResolveTenant(context.Background(), "shared-t")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info2.IsolateMode != router.IsolateModeShared {
		t.Fatalf("expected shared mode, got %s", info2.IsolateMode)
	}
}
