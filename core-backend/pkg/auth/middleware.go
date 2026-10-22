package auth

import (
	"errors"
	"strings"

	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/gofiber/fiber/v3"
)

// ExtractBearerToken parses the Bearer scheme from the Authorization header.
func ExtractBearerToken(c fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// JWTAuthMiddleware enforces authentication via JWT Bearer token.
func JWTAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenStr := ExtractBearerToken(c)

		if tokenStr != "" {
			claims, err := ValidateToken(jwtSecret, tokenStr)
			if err != nil {
				if errors.Is(err, ErrExpiredToken) {
					return exceptions.NewUnauthorizeError(response_code.AccessTokenExpired, nil)
				}
				return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
			}

			if claims.TokenType != TokenTypeAccess {
				return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
			}

			effectiveTenantID := claims.TenantID
			if claims.IsHost {
				if reqTenant := c.Get("X-Tenant-ID"); reqTenant != "" && reqTenant != "default" {
					effectiveTenantID = reqTenant
				}
			}

			// Populate fiber locals for downstream handlers and RBAC middleware
			c.Locals("user_id", claims.Subject)
			c.Locals("tenant_id", effectiveTenantID)
			c.Locals("email", claims.Email)
			c.Locals("user_role", claims.Role)
			c.Locals("is_host", claims.IsHost)
			c.Locals("token_claims", claims)

			// Propagate tenant context to Go std context
			if effectiveTenantID != "" {
				c.SetContext(tenantctx.WithContext(c.Context(), effectiveTenantID))
			}

			return c.Next()
		}

		// Fallback for developer DX / offline mock testing if enabled
		if rbac.IsLegacyHeaderAuthAllowed() {
			callerRole, _ := rbac.ExtractCallerContext(c)
			if callerRole != "" {
				return c.Next()
			}
		}

		return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
	}
}

// OptionalJWTAuthMiddleware inspects and sets JWT claims if a valid token is provided,
// but does not reject unauthenticated requests.
func OptionalJWTAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenStr := ExtractBearerToken(c)
		if tokenStr != "" {
			claims, err := ValidateToken(jwtSecret, tokenStr)
			if err == nil && claims.TokenType == TokenTypeAccess {
				effectiveTenantID := claims.TenantID
				if claims.IsHost {
					if reqTenant := c.Get("X-Tenant-ID"); reqTenant != "" && reqTenant != "default" {
						effectiveTenantID = reqTenant
					}
				}

				c.Locals("user_id", claims.Subject)
				c.Locals("tenant_id", effectiveTenantID)
				c.Locals("email", claims.Email)
				c.Locals("user_role", claims.Role)
				c.Locals("is_host", claims.IsHost)
				c.Locals("token_claims", claims)

				if effectiveTenantID != "" {
					c.SetContext(tenantctx.WithContext(c.Context(), effectiveTenantID))
				}
			}
		}
		return c.Next()
	}
}
