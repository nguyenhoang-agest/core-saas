package handler

import (
	"log/slog"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions"
	"github.com/example/core-backend/pkg/helpers"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/rbac"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service service.ITenantService
}

func NewHandler(service service.ITenantService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/", rbac.RequirePermission("tenant:list"), h.GetAllTenant)
	r.Get("/:id", rbac.RequirePermission("tenant:view"), h.GetByID)
	r.Post("/", rbac.RequirePermission("tenant:create"), h.Create)
	r.Put("/:id", rbac.RequirePermission("tenant:update"), h.Update)
	r.Patch("/:id/status", rbac.RequirePermission("tenant:update"), h.UpdateStatus)
	r.Delete("/:id", rbac.RequirePermission("tenant:delete"), h.Delete)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input dto.CreateTenantRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	tenant, err := h.service.Create(c.Context(), &input)
	if err != nil {
		slog.ErrorContext(c.Context(), "Error during create Tenant", "error", err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(httpx.NewResponse(tenant))
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	isHost, _ := c.Locals("is_host").(bool)
	callerTenantID, _ := c.Locals("tenant_id").(string)

	if !isHost && callerTenantID != "" && callerTenantID != id {
		return exceptions.NewForbiddenError(response_code.InsufficientPermission)
	}

	tenant, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(tenant))
}

func (h *Handler) GetAllTenant(c fiber.Ctx) error {
	var opts dto.PaginationOptions
	if err := c.Bind().Query(&opts); err != nil {
		slog.ErrorContext(c.Context(), "Query validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, opts)
	}

	result, err := h.service.GetAllTenant(c.Context(), &opts)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(result))
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	isHost, _ := c.Locals("is_host").(bool)
	callerTenantID, _ := c.Locals("tenant_id").(string)

	if !isHost && callerTenantID != "" && callerTenantID != id {
		return exceptions.NewForbiddenError(response_code.InsufficientPermission)
	}

	var input dto.UpdateTenantRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	tenant, err := h.service.Update(c.Context(), id, &input)
	if err != nil {
		slog.ErrorContext(c.Context(), "Error during update Tenant", "id", id, "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(tenant))
}

func (h *Handler) UpdateStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var input dto.UpdateTenantStatusRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	tenant, err := h.service.UpdateStatus(c.Context(), id, input.Status)
	if err != nil {
		slog.ErrorContext(c.Context(), "Error during update Tenant status", "id", id, "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(tenant))
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		slog.ErrorContext(c.Context(), "Error during delete Tenant", "id", id, "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(fiber.Map{
		"deleted": true,
		"id":      id,
	}))
}
