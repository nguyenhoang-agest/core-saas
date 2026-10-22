package tenantctx

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

// DefaultTenantID is used as fallback when no X-Tenant-ID header is present.
const DefaultTenantID = "default"

func FiberMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenantID := c.Get("x-tenant-id")
		if tenantID == "" {
			tenantID = DefaultTenantID
		}

		c.Locals("tenant_id", tenantID)
		c.SetContext(WithContext(c.Context(), tenantID))

		return c.Next()
	}
}

func GetTenantID(ctx context.Context) string {
	if id, ok := FromContext(ctx); ok && id != "" {
		return id
	}
	return DefaultTenantID
}
