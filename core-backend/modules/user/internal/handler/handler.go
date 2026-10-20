package handler

import (
	"log/slog"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/model"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/helpers"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service service.IUserService
}

func NewHandler(service service.IUserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/", h.GetAll)
	r.Get("/roles/meta", h.GetRolesMeta)
	r.Get("/permissions/meta", h.GetPermissionsMeta)
	r.Get("/:id", h.GetByID)

	// Protected mutating endpoints
	r.Post("/", rbac.RequirePermission("iam:users:manage"), h.Create)
	r.Put("/:id", rbac.RequirePermission("iam:users:manage"), h.Update)
	r.Patch("/:id/status", rbac.RequirePermission("iam:users:manage"), h.UpdateStatus)
	r.Delete("/:id", rbac.RequirePermission("iam:users:manage"), h.Delete)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input dto.CreateUserRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}
	entity, err := h.service.Create(c.Context(), &input)
	if err != nil {
		slog.ErrorContext(c.Context(), "Error during create user", "error", err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(httpx.NewResponse(entity))
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	entity, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(entity))
}

func (h *Handler) GetAll(c fiber.Ctx) error {
	var opts dto.PaginationOptions
	if err := c.Bind().Query(&opts); err != nil {
		slog.ErrorContext(c.Context(), "Query validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, opts)
	}

	result, err := h.service.GetAll(c.Context(), &opts)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(result))
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var input dto.UpdateUserRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	entity, err := h.service.Update(c.Context(), id, &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(entity))
}

func (h *Handler) UpdateStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var input dto.UpdateUserStatusRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	entity, err := h.service.UpdateStatus(c.Context(), id, model.Status(input.Status))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(entity))
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(map[string]bool{"deleted": true}))
}

func (h *Handler) GetRolesMeta(c fiber.Ctx) error {
	roles := h.service.GetRolesMeta()
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(roles))
}

func (h *Handler) GetPermissionsMeta(c fiber.Ctx) error {
	perms := h.service.GetPermissionsMeta()
	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(perms))
}
