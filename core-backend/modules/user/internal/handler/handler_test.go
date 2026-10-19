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
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/i18n"
	"github.com/gofiber/fiber/v3"
)

type mockUserService struct {
	createFunc  func(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
	getByIDFunc func(ctx context.Context, id string) (*model.User, error)
	getAllFunc  func(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.User], error)
	deleteFunc  func(ctx context.Context, id string) error
}

func (m *mockUserService) Create(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserService) GetAll(ctx context.Context, req *dto.PaginationOptions) (*httpx.PaginatedData[model.User], error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockUserService) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func setupTestApp(svc *mockUserService) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:    httpx.ErrorHandlerFnc,
		StructValidator: httpx.NewStructValidator(),
	})
	app.Use(i18n.LocaleMiddleware())
	h := handler.NewHandler(svc)
	h.RegisterRoutes(app)
	return app
}

func TestCreateUser_Success(t *testing.T) {
	now := time.Now().UTC()
	expected := &model.User{
		BaseModel: model.BaseModel{
			Id: "test-id-1",
			AuditModel: model.AuditModel{
				CreateAt: now,
				UpdateAt: now,
			},
		},
		Name:        "Test Item",
		Description: "A test description",
		Status:      model.StatusActive,
	}

	mockSvc := &mockUserService{
		createFunc: func(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
			if req.Name != "Test Item" {
				t.Fatalf("expected name 'Test Item', got '%s'", req.Name)
			}
			return expected, nil
		},
	}

	app := setupTestApp(mockSvc)

	body := `{"name":"Test Item","description":"A test description"}`
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

func TestCreateUser_ValidationError(t *testing.T) {
	mockSvc := &mockUserService{}
	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":""}`))
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

	if errResp.MessageCode != response_code.CommonError.Code() {
		t.Fatalf("expected code %s, got %s", response_code.CommonError.Code(), errResp.MessageCode)
	}
}

func TestGetByID_Success(t *testing.T) {
	expected := &model.User{
		BaseModel: model.BaseModel{
			Id: "item-123",
		},
		Name:   "Test Item",
		Status: model.StatusActive,
	}

	mockSvc := &mockUserService{
		getByIDFunc: func(ctx context.Context, id string) (*model.User, error) {
			if id != "item-123" {
				t.Fatalf("expected id 'item-123', got '%s'", id)
			}
			return expected, nil
		},
	}

	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/item-123", nil)
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
		getByIDFunc: func(ctx context.Context, id string) (*model.User, error) {
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
