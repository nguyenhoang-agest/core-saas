package rbac_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/gofiber/fiber/v3"
)

func TestPermissionRegistry_Validation(t *testing.T) {
	reg := rbac.GetRegistry()

	// 1. Valid concrete keys pass validation
	err := reg.ValidateConcreteKeys([]string{"iam:view", "billing:view", "billing:invoices:view"})
	if err != nil {
		t.Fatalf("expected valid keys to pass, got: %v", err)
	}

	// 2. Wildcards are rejected
	errWildcard := reg.ValidateConcreteKeys([]string{"billing:*"})
	if errWildcard == nil {
		t.Fatal("expected wildcard 'billing:*' to be rejected")
	}

	errStar := reg.ValidateConcreteKeys([]string{"*"})
	if errStar == nil {
		t.Fatal("expected universal wildcard '*' to be rejected")
	}

	// 3. Unknown keys are rejected
	errUnknown := reg.ValidateConcreteKeys([]string{"unknown:fake:perm"})
	if errUnknown == nil {
		t.Fatal("expected unknown key to be rejected")
	}
}

func TestScopeFilteredWildcards(t *testing.T) {
	reg := rbac.GetRegistry()

	// 1. Non-host tenant expansion: must NEVER include ScopePlatform keys
	nonHostPerms := reg.ExpandPermissions([]string{"*"}, false)
	for _, p := range nonHostPerms {
		if p == "tenant:create" || p == "tenant:delete" {
			t.Fatalf("scope boundary violated: non-host tenant expansion included platform permission %s", p)
		}
	}

	// Verify standard tenant permissions ARE present
	hasTenantView := false
	for _, p := range nonHostPerms {
		if p == "tenant:view" {
			hasTenantView = true
			break
		}
	}
	if !hasTenantView {
		t.Fatal("expected tenant:view to be included in non-host tenant expansion")
	}

	// 2. Host tenant expansion: includes ScopePlatform keys
	hostPerms := reg.ExpandPermissions([]string{"*"}, true)
	hasTenantDelete := false
	for _, p := range hostPerms {
		if p == "tenant:delete" {
			hasTenantDelete = true
			break
		}
	}
	if !hasTenantDelete {
		t.Fatal("expected host tenant expansion to include platform permission tenant:delete")
	}

	// 3. HasPermissionWithScope for Owner
	if rbac.HasPermissionWithScope("owner", nil, false, "tenant:delete") {
		t.Fatal("security violation: non-host owner must not have tenant:delete")
	}
	if !rbac.HasPermissionWithScope("owner", nil, true, "tenant:delete") {
		t.Fatal("expected host owner to have tenant:delete")
	}
	if !rbac.HasPermissionWithScope("owner", nil, false, "iam:users:manage") {
		t.Fatal("expected non-host owner to have tenant-scoped iam:users:manage")
	}
}

func TestRequirePermission_FailClosed_And_SpoofedHeaders(t *testing.T) {
	// Disable legacy header auth to enforce fail-closed security
	rbac.SetAllowLegacyHeaderAuth(false)
	defer rbac.SetAllowLegacyHeaderAuth(true) // reset for other tests

	app := fiber.New(fiber.Config{
		ErrorHandler: httpx.ErrorHandlerFnc,
	})

	app.Use(rbac.StripExternalAuthHeaders())

	app.Get("/secure-endpoint", rbac.RequirePermission("iam:roles:manage"), func(c fiber.Ctx) error {
		return c.SendString("success")
	})

	// Test 1: Spoofed X-User-Role: owner with no verified user_id -> 401 Unauthorized
	req1 := httptest.NewRequest(http.MethodGet, "/secure-endpoint", nil)
	req1.Header.Set("X-User-Role", "owner")
	req1.Header.Set("X-User-Permissions", "*")
	resp1, err := app.Test(req1)
	if err != nil || resp1.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated spoofed request, got %d, err %v", resp1.StatusCode, err)
	}

	// Test 2: Spoofed headers stripped and rejected
	req2 := httptest.NewRequest(http.MethodGet, "/secure-endpoint", nil)
	req2.Header.Set("X-User-Permissions", "iam:roles:manage")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for spoofed permissions header, got %d", resp2.StatusCode)
	}
}

type mockResolver struct {
	allowed map[string]bool
}

func (m *mockResolver) ResolvePermissions(ctx context.Context, tenantID, userID string, isHost bool) ([]string, error) {
	var list []string
	for k := range m.allowed {
		list = append(list, k)
	}
	return list, nil
}

func (m *mockResolver) HasPermission(ctx context.Context, tenantID, userID string, isHost bool, requiredKey string) (bool, error) {
	return m.allowed[requiredKey], nil
}

func (m *mockResolver) Invalidate(ctx context.Context, tenantID, userID string) {}

func TestRequirePermission_WithResolver(t *testing.T) {
	resolver := &mockResolver{
		allowed: map[string]bool{
			"billing:view": true,
		},
	}
	rbac.SetPermissionResolver(resolver)
	defer rbac.SetPermissionResolver(nil)

	app := fiber.New(fiber.Config{
		ErrorHandler: httpx.ErrorHandlerFnc,
	})

	// Simulate verified auth middleware injecting locals
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", "user-123")
		c.Locals("tenant_id", "tenant-alpha")
		c.Locals("is_host", false)
		return c.Next()
	})

	app.Get("/billing", rbac.RequirePermission("billing:view"), func(c fiber.Ctx) error {
		return c.SendString("billing ok")
	})
	app.Get("/admin", rbac.RequirePermission("iam:roles:manage"), func(c fiber.Ctx) error {
		return c.SendString("admin ok")
	})

	// Billing should be granted
	req1 := httptest.NewRequest(http.MethodGet, "/billing", nil)
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from resolver, got %d", resp1.StatusCode)
	}

	// Admin should be forbidden
	req2 := httptest.NewRequest(http.MethodGet, "/admin", nil)
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 from resolver for unheld permission, got %d", resp2.StatusCode)
	}
}
