package tenantctx

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

// DefaultTenantID is used as fallback when no X-Tenant-ID header is present.
const DefaultTenantID = "default"

// FiberMiddleware extracts the active tenant ID from the X-Tenant-ID HTTP header,
// stores it in Fiber context locals, and wraps it into the Go request context.
func FiberMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = DefaultTenantID
		}

		c.Locals("tenant_id", tenantID)
		c.SetContext(WithContext(c.Context(), tenantID))

		return c.Next()
	}
}

// GetTenantID retrieves the tenant ID from context, or returns DefaultTenantID if not set.
func GetTenantID(ctx context.Context) string {
	if id, ok := FromContext(ctx); ok && id != "" {
		return id
	}
	return DefaultTenantID
}
