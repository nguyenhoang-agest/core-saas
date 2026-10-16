package i18n

import (
	"github.com/gofiber/fiber/v3"
)

// LocaleMiddleware resolves client locale from Accept-Language, X-Locale, or ?lang= query param.
// It stores the normalized locale ("en" or "vi") in Fiber locals and request Context.
func LocaleMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		raw := c.Query("lang")
		if raw == "" {
			raw = c.Get("X-Locale")
		}
		if raw == "" {
			raw = c.Get("Accept-Language")
		}

		locale := NormalizeLocale(raw)

		c.Locals("locale", locale)
		c.SetContext(WithLocale(c.Context(), locale))

		return c.Next()
	}
}
