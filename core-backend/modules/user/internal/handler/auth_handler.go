package handler

import (
	"log/slog"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/auth"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/helpers"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService service.IAuthService
	jwtSecret   string
}

func NewAuthHandler(authService service.IAuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}

func (h *AuthHandler) RegisterRoutes(r fiber.Router) {
	group := r.Group("/auth")

	// Public auth endpoints
	group.Post("/register", h.RegisterOrg)
	group.Post("/login", h.Login)
	group.Post("/refresh", h.Refresh)

	// Protected auth endpoints
	protected := group.Group("", auth.JWTAuthMiddleware(h.jwtSecret))
	protected.Get("/me", h.Me)
	protected.Post("/logout", h.Logout)
}

func (h *AuthHandler) RegisterOrg(c fiber.Ctx) error {
	var input dto.RegisterOrgRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "body validation error in register", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	result, err := h.authService.RegisterOrg(c.Context(), &input)
	if err != nil {
		slog.ErrorContext(c.Context(), "error during register org", "error", err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(httpx.NewResponse(result))
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var input dto.LoginRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "body validation error in login", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	result, err := h.authService.Login(c.Context(), &input)
	if err != nil {
		slog.WarnContext(c.Context(), "login attempt failed", "email", input.Email, "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(result))
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var input dto.RefreshTokenRequest
	if err := c.Bind().Body(&input); err != nil {
		return helpers.ValidationErrorFromBind(err, input)
	}

	result, err := h.authService.RefreshToken(c.Context(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(result))
}

func (h *AuthHandler) Me(c fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	tenantID, _ := c.Locals("tenant_id").(string)
	isHost, _ := c.Locals("is_host").(bool)

	if userID == "" || tenantID == "" {
		return exceptions.NewUnauthorizeError(response_code.InvalidToken, nil)
	}

	result, err := h.authService.GetCurrentUser(c.Context(), tenantID, userID, isHost)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(result))
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(fiber.Map{
		"message": "successfully logged out",
	}))
}
