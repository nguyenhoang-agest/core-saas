package rbac_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/gofiber/fiber/v3"
)

func TestRequirePermission_Middleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: httpx.ErrorHandlerFnc,
	})

	app.Get("/protected", rbac.RequirePermission("iam:users:manage"), func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Test 1: Caller with default owner role (no headers) -> Allowed
	req1 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp1, err := app.Test(req1)
	if err != nil || resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for default owner, got %d, err %v", resp1.StatusCode, err)
	}

	// Test 2: Caller with member role and no custom ACL -> Forbidden 403
	req2 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req2.Header.Set("X-User-Role", "member")
	resp2, err := app.Test(req2)
	if err != nil || resp2.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for member without ACL, got %d", resp2.StatusCode)
	}

	// Test 3: Caller with member role AND custom ACL permission -> Allowed 200
	req3 := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req3.Header.Set("X-User-Role", "member")
	req3.Header.Set("X-User-Permissions", "iam:users:manage,billing:view")
	resp3, err := app.Test(req3)
	if err != nil || resp3.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for member with ACL, got %d", resp3.StatusCode)
	}
}

func TestRequireRole_Middleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: httpx.ErrorHandlerFnc,
	})

	app.Get("/admin-only", rbac.RequireRole(string(rbac.RoleAdmin)), func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Admin -> Allowed
	req1 := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req1.Header.Set("X-User-Role", "admin")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for admin, got %d", resp1.StatusCode)
	}

	// Viewer -> Forbidden 403
	req2 := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req2.Header.Set("X-User-Role", "viewer")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer, got %d", resp2.StatusCode)
	}
}
