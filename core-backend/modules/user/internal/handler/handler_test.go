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

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/errcode"
	"github.com/example/core-backend/modules/user/internal/handler"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/i18n"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/gofiber/fiber/v3"
)

type mockUserService struct {
	createFunc       func(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	getByIDFunc      func(ctx context.Context, id string) (*dto.UserResponse, error)
	getAllFunc       func(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[dto.UserResponse], error)
	updateFunc       func(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	updateStatusFunc func(ctx context.Context, id string, status model.Status) (*dto.UserResponse, error)
	deleteFunc       func(ctx context.Context, id string) error
}

func (m *mockUserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUserService) GetByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserService) GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[dto.UserResponse], error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUserService) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockUserService) UpdateStatus(ctx context.Context, id string, status model.Status) (*dto.UserResponse, error) {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status)
	}
	return nil, nil
}

func (m *mockUserService) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockUserService) GetRolesMeta() []rbac.RoleDefinition {
	return rbac.GetSystemRoles()
}

func (m *mockUserService) GetPermissionsMeta() []rbac.PermissionDefinition {
	return rbac.SystemPermissionsCatalog
}

func setupTestApp(svc *mockUserService) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:    httpx.ErrorHandlerFnc,
		StructValidator: httpx.NewStructValidator(),
	})
	app.Use(tenantctx.FiberMiddleware())
	app.Use(i18n.LocaleMiddleware())
	h := handler.NewHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestCreateUser_Success(t *testing.T) {
	now := time.Now().UTC()
	expected := &dto.UserResponse{
		Id:          "test-id-1",
		TenantID:    "default",
		Email:       "alice@example.com",
		Name:        "Alice User",
		Role:        "admin",
		Permissions: []string{"iam:users:manage"},
		Status:      model.StatusActive,
		CreateAt:    now,
		UpdateAt:    now,
	}

	mockSvc := &mockUserService{
		createFunc: func(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
			if req.Email != "alice@example.com" {
				t.Fatalf("expected email 'alice@example.com', got '%s'", req.Email)
			}
			return expected, nil
		},
	}

	app := setupTestApp(mockSvc)

	body := `{"email":"alice@example.com","name":"Alice User","role":"admin"}`
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

func TestCreateUser_Forbidden(t *testing.T) {
	mockSvc := &mockUserService{}
	app := setupTestApp(mockSvc)

	body := `{"email":"alice@example.com","name":"Alice User","role":"member"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	// Caller is a viewer without iam:users:manage
	req.Header.Set("X-User-Role", "viewer")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error making test request: %v", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden, got %d", resp.StatusCode)
	}
}

func TestGetByID_Success(t *testing.T) {
	expected := &dto.UserResponse{
		Id:     "user-123",
		Email:  "user@example.com",
		Name:   "Test User",
		Role:   "member",
		Status: model.StatusActive,
	}

	mockSvc := &mockUserService{
		getByIDFunc: func(ctx context.Context, id string) (*dto.UserResponse, error) {
			if id != "user-123" {
				t.Fatalf("expected id 'user-123', got '%s'", id)
			}
			return expected, nil
		},
	}

	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user-123", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	mockSvc := &mockUserService{
		getByIDFunc: func(ctx context.Context, id string) (*dto.UserResponse, error) {
			return nil, exceptions.NewNotFoundError(errcode.NotFoundUser)
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
}

func TestGetRolesMeta(t *testing.T) {
	mockSvc := &mockUserService{}
	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/roles/meta", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", resp.StatusCode)
	}
}
