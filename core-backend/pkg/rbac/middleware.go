package rbac

import (
	"strings"
	"sync"

	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/gofiber/fiber/v3"
)

var (
	legacyHeaderMu        sync.RWMutex
	allowLegacyHeaderAuth = true // Default true to preserve developer DX / offline tests; gated in production configs
)

// SetAllowLegacyHeaderAuth toggles whether unverified request headers (X-User-Role, X-User-Permissions)
// are accepted. Must be set to false in production environments to fail closed against spoofed headers.
func SetAllowLegacyHeaderAuth(allowed bool) {
	legacyHeaderMu.Lock()
	defer legacyHeaderMu.Unlock()
	allowLegacyHeaderAuth = allowed
}

// IsLegacyHeaderAuthAllowed checks whether legacy header-based caller context is permitted.
func IsLegacyHeaderAuthAllowed() bool {
	legacyHeaderMu.RLock()
	defer legacyHeaderMu.RUnlock()
	return allowLegacyHeaderAuth
}

// StripExternalAuthHeaders returns a Fiber middleware that strips potentially spoofed auth headers
// unless legacy header authentication is explicitly permitted by configuration.
func StripExternalAuthHeaders() fiber.Handler {
	return func(c fiber.Ctx) error {
		if !IsLegacyHeaderAuthAllowed() {
			c.Request().Header.Del("X-User-Role")
			c.Request().Header.Del("X-User-Permissions")
			c.Request().Header.Del("X-User-ID")
		}
		return c.Next()
	}
}

// ExtractCallerContext extracts user role and ACL permissions from the Fiber context or request headers.
// If legacy header authentication is disabled, it only reads from verified context locals.
func ExtractCallerContext(c fiber.Ctx) (role string, permissions []string) {
	if r, ok := c.Locals("user_role").(string); ok && r != "" {
		role = r
	} else if IsLegacyHeaderAuthAllowed() {
		if r := c.Get("X-User-Role"); r != "" {
			role = r
		} else {
			role = string(RoleOwner)
		}
	}

	if perms, ok := c.Locals("user_permissions").([]string); ok {
		permissions = perms
	} else if IsLegacyHeaderAuthAllowed() {
		if raw := c.Get("X-User-Permissions"); raw != "" {
			parts := strings.Split(raw, ",")
			for _, p := range parts {
				trimmed := strings.TrimSpace(p)
				if trimmed != "" {
					permissions = append(permissions, trimmed)
				}
			}
		}
	}

	return role, permissions
}

// RequirePermission constructs a Fiber v3 handler that verifies whether the caller
// possesses the specified permission. Fails closed if caller identity is unverified
// and legacy header fallback is disabled.
func RequirePermission(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		resolver := GetPermissionResolver()
		userID, _ := c.Locals("user_id").(string)
		tenantID, _ := c.Locals("tenant_id").(string)
		isHost, _ := c.Locals("is_host").(bool)

		// 1. Primary path: Evaluate dynamically via registered PermissionResolver
		if resolver != nil && userID != "" && tenantID != "" {
			hasPerm, err := resolver.HasPermission(c.Context(), tenantID, userID, isHost, permission)
			if err != nil {
				return exceptions.NewInternalServerError(response_code.InternalServerGeneral)
			}
			if !hasPerm {
				return exceptions.NewForbiddenError(response_code.InsufficientPermission)
			}
			return c.Next()
		}

		// 2. Fail-closed security boundary: If legacy header auth is disabled, deny unverified requests
		if !IsLegacyHeaderAuthAllowed() {
			if userID == "" {
				return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
			}
			return exceptions.NewForbiddenError(response_code.InsufficientPermission)
		}

		// 3. Fallback path for dev/test environments with AllowLegacyHeaderAuth = true
		role, perms := ExtractCallerContext(c)
		if !HasPermission(role, perms, permission) {
			return exceptions.NewForbiddenError(response_code.InsufficientPermission)
		}

		return c.Next()
	}
}

// RequireRole constructs a Fiber v3 handler that verifies whether the caller
// is assigned at least one of the specified roles (or is the workspace Owner).
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, _ := c.Locals("user_id").(string)

		if !IsLegacyHeaderAuthAllowed() && userID == "" {
			return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
		}

		role, _ := ExtractCallerContext(c)

		if role == string(RoleOwner) {
			return c.Next()
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(role, allowed) {
				return c.Next()
			}
		}

		return exceptions.NewForbiddenError(response_code.AccessDenied)
	}
}
