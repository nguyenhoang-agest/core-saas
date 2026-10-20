package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/handler"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/i18n"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/gofiber/fiber/v3"
)

type mockRoleService struct {
	listRolesFunc             func(ctx context.Context) ([]dto.RoleResponse, error)
	getRoleByIDFunc           func(ctx context.Context, id string) (*dto.RoleResponse, error)
	createRoleFunc            func(ctx context.Context, actorPerms []string, isHost bool, req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	updateRoleFunc            func(ctx context.Context, actorPerms []string, isHost bool, id string, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	deleteRoleFunc            func(ctx context.Context, id string) (int, error)
	cloneRoleFunc             func(ctx context.Context, actorPerms []string, isHost bool, id string, req *dto.CloneRoleRequest) (*dto.RoleResponse, error)
	assignUserRolesFunc       func(ctx context.Context, actorRole string, actorPerms []string, isHost bool, targetUserID string, req *dto.AssignUserRolesRequest, assignedBy string) error
	getPermissionsCatalogFunc func(ctx context.Context, isHost bool) []dto.PermissionResponse
}

func (m *mockRoleService) ListRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	if m.listRolesFunc != nil {
		return m.listRolesFunc(ctx)
	}
	return []dto.RoleResponse{}, nil
}

func (m *mockRoleService) GetRoleByID(ctx context.Context, id string) (*dto.RoleResponse, error) {
	if m.getRoleByIDFunc != nil {
		return m.getRoleByIDFunc(ctx, id)
	}
	return &dto.RoleResponse{ID: id, Name: "Test Role"}, nil
}

func (m *mockRoleService) CreateRole(ctx context.Context, actorPerms []string, isHost bool, req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	if m.createRoleFunc != nil {
		return m.createRoleFunc(ctx, actorPerms, isHost, req)
	}
	return &dto.RoleResponse{
		ID:          "new-role-id",
		Slug:        req.Slug,
		Name:        req.Name,
		Permissions: req.Permissions,
		Version:     1,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}, nil
}

func (m *mockRoleService) UpdateRole(ctx context.Context, actorPerms []string, isHost bool, id string, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	if m.updateRoleFunc != nil {
		return m.updateRoleFunc(ctx, actorPerms, isHost, id, req)
	}
	return &dto.RoleResponse{ID: id, Name: req.Name, Version: req.Version + 1}, nil
}

func (m *mockRoleService) DeleteRole(ctx context.Context, id string) (int, error) {
	if m.deleteRoleFunc != nil {
		return m.deleteRoleFunc(ctx, id)
	}
	return 0, nil
}

func (m *mockRoleService) CloneRole(ctx context.Context, actorPerms []string, isHost bool, id string, req *dto.CloneRoleRequest) (*dto.RoleResponse, error) {
	if m.cloneRoleFunc != nil {
		return m.cloneRoleFunc(ctx, actorPerms, isHost, id, req)
	}
	return &dto.RoleResponse{ID: "cloned-id", Slug: req.Slug, Name: req.Name}, nil
}

func (m *mockRoleService) AssignUserRoles(ctx context.Context, actorRole string, actorPerms []string, isHost bool, targetUserID string, req *dto.AssignUserRolesRequest, assignedBy string) error {
	if m.assignUserRolesFunc != nil {
		return m.assignUserRolesFunc(ctx, actorRole, actorPerms, isHost, targetUserID, req, assignedBy)
	}
	return nil
}

func (m *mockRoleService) GetPermissionsCatalog(ctx context.Context, isHost bool) []dto.PermissionResponse {
	if m.getPermissionsCatalogFunc != nil {
		return m.getPermissionsCatalogFunc(ctx, isHost)
	}
	return []dto.PermissionResponse{
		{Key: "iam:view", Module: "iam", Scope: "tenant"},
		{Key: "billing:view", Module: "billing", Scope: "tenant"},
	}
}

func setupRoleTestApp(svc *mockRoleService) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:    httpx.ErrorHandlerFnc,
		StructValidator: httpx.NewStructValidator(),
	})
	app.Use(tenantctx.FiberMiddleware())
	app.Use(i18n.LocaleMiddleware())

	h := handler.NewRoleHandler(svc)
	group := app.Group("/users")
	h.RegisterRoutes(group)

	return app
}

func TestRoleHandler_ListAndPermissions(t *testing.T) {
	svc := &mockRoleService{}
	app := setupRoleTestApp(svc)

	// GET /users/permissions
	req1 := httptest.NewRequest(http.MethodGet, "/users/permissions", nil)
	req1.Header.Set("X-User-Role", "admin")
	resp1, err := app.Test(req1)
	if err != nil || resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for permissions catalog, got %d, err %v", resp1.StatusCode, err)
	}

	// GET /users/roles
	req2 := httptest.NewRequest(http.MethodGet, "/users/roles", nil)
	req2.Header.Set("X-User-Role", "admin")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for roles list, got %d", resp2.StatusCode)
	}
}

func TestRoleHandler_CreateRole(t *testing.T) {
	svc := &mockRoleService{}
	app := setupRoleTestApp(svc)

	payload := dto.CreateRoleRequest{
		Name:        "Auditor",
		Slug:        "auditor",
		Permissions: []string{"iam:view", "billing:view"},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/users/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "admin")

	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created for role creation, got %d, err %v", resp.StatusCode, err)
	}
}

func TestRoleHandler_AssignRoles(t *testing.T) {
	svc := &mockRoleService{}
	app := setupRoleTestApp(svc)

	payload := dto.AssignUserRolesRequest{
		RoleIDs: []string{"role-1", "role-2"},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/users/target-user-123/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "admin")

	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for user role assignment, got %d, err %v", resp.StatusCode, err)
	}
}
