package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/errcode"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/exceptions"
)

type mockRepository struct {
	createdTenant   *model.Tenant
	getByNameResult *model.Tenant
	getByIDResult   *model.Tenant
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
}
