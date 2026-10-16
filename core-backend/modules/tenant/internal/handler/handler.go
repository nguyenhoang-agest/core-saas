package handler

import (
	"log/slog"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/helpers"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service service.ITenantService
}

func NewHandler(service service.ITenantService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.GetAllTenant)
	r.Get("/:id", h.GetByID)
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
