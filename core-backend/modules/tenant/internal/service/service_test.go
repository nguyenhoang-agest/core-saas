package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/errcode"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/db"
	"github.com/example/core-backend/pkg/exceptions"
)

type mockRepository struct {
	createdTenant   *model.Tenant
	getByNameResult *model.Tenant
	getByIDResult   *model.Tenant
	listResult      []model.Tenant
	listTotal       int64
}

func (m *mockRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	m.createdTenant = tenant
	return nil
}

func (m *mockRepository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	return m.getByIDResult, nil
}

func (m *mockRepository) GetTenantByName(ctx context.Context, name string) (*model.Tenant, error) {
	return m.getByNameResult, nil
}

func (m *mockRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockRepository) GetAllTenant(ctx context.Context, req *dto.PaginationOptions) ([]model.Tenant, int64, error) {
	if m.listResult != nil {
		return m.listResult, m.listTotal, nil
	}
	return []model.Tenant{}, 0, nil
}

func TestService_Create_DefaultShared(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	req := &dto.CreateTenantRequest{
		Name:        "Test Shared",
		Description: "Shared Org",
	}

	tenant, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tenant.IsolateMode != model.IsolateModeShared {
		t.Fatalf("expected isolate mode 'shared', got '%s'", tenant.IsolateMode)
	}
}

func TestService_Create_InvalidIsolateMode(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	req := &dto.CreateTenantRequest{
		Name:        "Test Invalid",
		IsolateMode: "unknown_mode",
	}

	_, err := svc.Create(context.Background(), req)
	var badReqErr exceptions.BadRequestError
	if !errors.As(err, &badReqErr) || badReqErr.MessageCode != errcode.InvalidIsolateMode.Code() {
		t.Fatalf("expected BadRequestError with InvalidIsolateMode (%s), got %v", errcode.InvalidIsolateMode.Code(), err)
	}
}

func TestService_Create_Dedicated_MissingDSN(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	req := &dto.CreateTenantRequest{
		Name:        "Test Dedicated Fail",
		IsolateMode: model.IsolateModeDedicated,
	}

	_, err := svc.Create(context.Background(), req)
	var badReqErr exceptions.BadRequestError
	if !errors.As(err, &badReqErr) || badReqErr.MessageCode != errcode.MissingDedicatedDBConfig.Code() {
		t.Fatalf("expected MissingDedicatedDBConfig (%s) when DSN is missing, got %v", errcode.MissingDedicatedDBConfig.Code(), err)
	}

	// Also test empty DSN string
	req.DBConfig = &model.TenantDBConfig{DSN: "  "}
	_, err = svc.Create(context.Background(), req)
	if !errors.As(err, &badReqErr) || badReqErr.MessageCode != errcode.MissingDedicatedDBConfig.Code() {
		t.Fatalf("expected MissingDedicatedDBConfig (%s) when DSN is empty whitespace, got %v", errcode.MissingDedicatedDBConfig.Code(), err)
	}
}

func TestService_Create_Dedicated_Success(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	req := &dto.CreateTenantRequest{
		Name:        "Test Dedicated Success",
		IsolateMode: model.IsolateModeDedicated,
		DBConfig: &model.TenantDBConfig{
			DSN: "postgres://tenant_db:5432/app",
		},
	}

	tenant, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tenant.IsolateMode != model.IsolateModeDedicated {
		t.Fatalf("expected isolate mode 'dedicated', got '%s'", tenant.IsolateMode)
	}

	if tenant.DBConfig == nil || tenant.DBConfig.DSN != "postgres://tenant_db:5432/app" {
		t.Fatalf("unexpected DBConfig: %+v", tenant.DBConfig)
	}

	if tenant.DBConfig.MaxConns != 10 {
		t.Fatalf("expected default MaxConns to be 10, got %d", tenant.DBConfig.MaxConns)
	}

	if tenant.DBConfig.SSLMode != db.DefaultTenantSSLMode {
		t.Fatalf("expected default SSLMode to be %q, got %s", db.DefaultTenantSSLMode, tenant.DBConfig.SSLMode)
	}
}

func TestService_Create_Dedicated_DenyUserChoiceOfSSLMode(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	// Caller attempts to set SSLMode to "disable"
	req := &dto.CreateTenantRequest{
		Name:        "Test Insecure SSL",
		IsolateMode: model.IsolateModeDedicated,
		DBConfig: &model.TenantDBConfig{
			DSN:     "postgres://tenant_db:5432/app",
			SSLMode: "disable",
		},
	}

	tenant, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Must be forced to system-managed DefaultTenantSSLMode ("require")
	if tenant.DBConfig.SSLMode != db.DefaultTenantSSLMode {
		t.Fatalf("expected SSLMode to be system-managed %q, but got %q", db.DefaultTenantSSLMode, tenant.DBConfig.SSLMode)
	}
}

func TestService_Create_Dedicated_DenyUserChoiceOfConnectionPool(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	// Caller attempts to set custom MaxConns (e.g. 100)
	req := &dto.CreateTenantRequest{
		Name:        "Test Custom Pool",
		IsolateMode: model.IsolateModeDedicated,
		DBConfig: &model.TenantDBConfig{
			DSN:      "postgres://tenant_db:5432/app",
			MaxConns: 100,
		},
	}

	tenant, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Must be forced to system-managed DefaultTenantMaxConns (10)
	if tenant.DBConfig.MaxConns != db.DefaultTenantMaxConns {
		t.Fatalf("expected MaxConns to be system-managed %d, but got %d", db.DefaultTenantMaxConns, tenant.DBConfig.MaxConns)
	}
}

func TestService_Create_Dedicated_SSRFBlocked(t *testing.T) {
	repo := &mockRepository{}
	svc := service.NewTenantService(repo, nil)

	blockedHosts := []string{
		"postgres://user:pass@localhost:5432/db",
		"postgres://user:pass@127.0.0.1:5432/db",
		"postgres://user:pass@127.0.0.2:5432/db",
		"postgres://user:pass@169.254.169.254:5432/db",
	}

	for _, dsn := range blockedHosts {
		req := &dto.CreateTenantRequest{
			Name:        "Test SSRF",
			IsolateMode: model.IsolateModeDedicated,
			DBConfig: &model.TenantDBConfig{
				DSN: dsn,
			},
		}

		_, err := svc.Create(context.Background(), req)
		var badReqErr exceptions.BadRequestError
		if !errors.As(err, &badReqErr) || badReqErr.MessageCode != errcode.InvalidDedicatedDBConfig.Code() {
			t.Fatalf("expected InvalidDedicatedDBConfig for prohibited DSN %q, got %v", dsn, err)
		}
	}
}

func TestService_List_Success(t *testing.T) {
	expectedTenants := []model.Tenant{
		{
			BaseModel: model.BaseModel{Id: "t-1"},
			Name:      "Acme Corp",
			Status:    model.StatusActive,
		},
		{
			BaseModel: model.BaseModel{Id: "t-2"},
			Name:      "Beta Ltd",
			Status:    model.StatusActive,
		},
	}

	repo := &mockRepository{
		listResult: expectedTenants,
		listTotal:  2,
	}
	svc := service.NewTenantService(repo, nil)

	req := &dto.PaginationOptions{}
	result, err := svc.GetAllTenant(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result.Items))
	}

	if result.Pagination.TotalItems != 2 {
		t.Errorf("expected total_items 2, got %d", result.Pagination.TotalItems)
	}
	if result.Pagination.Page != 1 {
		t.Errorf("expected default page 1, got %d", result.Pagination.Page)
	}
	if result.Pagination.PageSize != 20 {
		t.Errorf("expected default page_size 20, got %d", result.Pagination.PageSize)
	}
}
