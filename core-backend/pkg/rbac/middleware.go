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

func SetAllowLegacyHeaderAuth(allowed bool) {
	legacyHeaderMu.Lock()
	defer legacyHeaderMu.Unlock()
	allowLegacyHeaderAuth = allowed
}

func IsLegacyHeaderAuthAllowed() bool {
	legacyHeaderMu.RLock()
	defer legacyHeaderMu.RUnlock()
	return allowLegacyHeaderAuth
}

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

func RequirePermission(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		resolver := GetPermissionResolver()
		userID, _ := c.Locals("user_id").(string)
		tenantID, _ := c.Locals("tenant_id").(string)
		isHost, _ := c.Locals("is_host").(bool)

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

		if !IsLegacyHeaderAuthAllowed() {
			if userID == "" {
				return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
			}
			return exceptions.NewForbiddenError(response_code.InsufficientPermission)
		}

		role, perms := ExtractCallerContext(c)
		if !HasPermission(role, perms, permission) {
			return exceptions.NewForbiddenError(response_code.InsufficientPermission)
		}

		return c.Next()
	}
}

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
