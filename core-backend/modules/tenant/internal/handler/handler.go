package handler

import (
	"errors"
	"strings"

	"github.com/example/core-backend/modules/tenant/internal/dto"
	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/constants/response_code"
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
	r.Get("/:id", h.GetByID)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input dto.CreateTenantRequest
	if err := c.Bind().Body(&input); err != nil {
		return h.handleError(c, response_code.InvalidPayload)
	}

	if strings.TrimSpace(input.Name) == "" {
		return h.handleError(c, response_code.TenantNameRequired)
	}

	tenant, err := h.service.Create(c.Context(), &input)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(httpx.NewResponse(tenant))
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	tenant, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(tenant))
}

func (h *Handler) handleError(c fiber.Ctx, err error) error {
	var respCode response_code.IResponseCode
	if errors.As(err, &respCode) {
		return c.Status(respCode.HTTPStatus()).JSON(httpx.NewResponseError(respCode.Code(), nil))
	}

	return c.Status(response_code.InternalServerGeneral.HTTPStatus()).JSON(
		httpx.NewResponseError(response_code.InternalServerGeneral.Code(), nil),
	)
}

