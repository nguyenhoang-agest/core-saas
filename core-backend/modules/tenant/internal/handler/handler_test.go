package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/handler"
	"github.com/example/core-backend/modules/tenant/internal/model"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/gofiber/fiber/v3"
)

type mockTenantService struct {
	createFunc  func(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error)
	getByIDFunc func(ctx context.Context, id string) (*model.Tenant, error)
}

func (m *mockTenantService) Create(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockTenantService) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func setupTestApp(svc *mockTenantService) *fiber.App {
	app := fiber.New()
	h := handler.NewHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestCreateTenant_Success(t *testing.T) {
	now := time.Now().UTC()
	expectedTenant := &model.Tenant{
		BaseModel: model.BaseModel{
			Id: "test-tenant-1",
			AuditModel: model.AuditModel{
				CreateAt: now,
				UpdateAt: now,
			},
		},
		Name:        "Acme Corp",
		Description: "A test organization",
		Status:      model.StatusActive,
	}

	mockSvc := &mockTenantService{
		createFunc: func(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error) {
			if req.Name != "Acme Corp" {
				t.Fatalf("expected name 'Acme Corp', got '%s'", req.Name)
			}
			return expectedTenant, nil
		},
	}

	app := setupTestApp(mockSvc)

	body := `{"name":"Acme Corp","description":"A test organization"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error making test request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result httpx.Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if result.Status != "success" {
		t.Fatalf("expected status 'success', got '%s'", result.Status)
	}
}

func TestCreateTenant_InvalidJSON(t *testing.T) {
	mockSvc := &mockTenantService{}
	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var errResp httpx.ResponseError
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if errResp.MessageCode != response_code.InvalidPayload.Code() {
		t.Fatalf("expected code %s, got %s", response_code.InvalidPayload.Code(), errResp.MessageCode)
	}
}

func TestCreateTenant_EmptyName(t *testing.T) {
	mockSvc := &mockTenantService{}
	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"   "}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var errResp httpx.ResponseError
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if errResp.MessageCode != response_code.TenantNameRequired.Code() {
		t.Fatalf("expected code %s, got %s", response_code.TenantNameRequired.Code(), errResp.MessageCode)
	}
}

func TestGetByID_Success(t *testing.T) {
	expectedTenant := &model.Tenant{
		BaseModel: model.BaseModel{
			Id: "tenant-123",
		},
		Name:   "Acme Corp",
		Status: model.StatusActive,
	}

	mockSvc := &mockTenantService{
		getByIDFunc: func(ctx context.Context, id string) (*model.Tenant, error) {
			if id != "tenant-123" {
				t.Fatalf("expected id 'tenant-123', got '%s'", id)
			}
			return expectedTenant, nil
		},
	}

	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/tenant-123", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	mockSvc := &mockTenantService{
		getByIDFunc: func(ctx context.Context, id string) (*model.Tenant, error) {
			return nil, model.ErrNotFound
		},
	}

	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/non-existent", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 Not Found, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var errResp httpx.ResponseError
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if errResp.MessageCode != response_code.NotFoundTenant.Code() {
		t.Fatalf("expected code %s, got %s", response_code.NotFoundTenant.Code(), errResp.MessageCode)
	}
}

func TestCreateTenant_DedicatedMode(t *testing.T) {
	now := time.Now().UTC()
	expectedTenant := &model.Tenant{
		BaseModel: model.BaseModel{
			Id: "test-tenant-dedicated",
			AuditModel: model.AuditModel{
				CreateAt: now,
				UpdateAt: now,
			},
		},
		Name:        "Dedicated Corp",
		Description: "An isolated organization",
		Status:      model.StatusActive,
		IsolateMode: model.IsolateModeDedicated,
		DBConfig: &model.TenantDBConfig{
			DSN:      "postgres://user:pass@host:5432/dedicated_db",
			MaxConns: 20,
		},
	}

	mockSvc := &mockTenantService{
		createFunc: func(ctx context.Context, req *dto.CreateTenantRequest) (*model.Tenant, error) {
			if req.IsolateMode != model.IsolateModeDedicated {
				t.Fatalf("expected isolate_mode dedicated, got '%s'", req.IsolateMode)
			}
			if req.DBConfig == nil || req.DBConfig.DSN != "postgres://user:pass@host:5432/dedicated_db" {
				t.Fatalf("unexpected db config: %+v", req.DBConfig)
			}
			return expectedTenant, nil
		},
	}

	app := setupTestApp(mockSvc)

	body := `{
		"name": "Dedicated Corp",
		"description": "An isolated organization",
		"isolate_mode": "dedicated",
		"db_config": {
			"dsn": "postgres://user:pass@host:5432/dedicated_db",
			"max_conns": 20
		}
	}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error making test request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d", resp.StatusCode)
	}
}

