package handler

import (
	"log/slog"

	"github.com/example/core-backend/modules/user/internal/dto"
	"github.com/example/core-backend/modules/user/internal/service"
	"github.com/example/core-backend/pkg/helpers"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service service.IUserService
}

func NewHandler(service service.IUserService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Post("/", h.Create)
	r.Get("/", h.GetAll)
	r.Get("/:id", h.GetByID)
	r.Delete("/:id", h.Delete)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var input dto.CreateUserRequest
	if err := c.Bind().Body(&input); err != nil {
		slog.ErrorContext(c.Context(), "Body validation error", "error", err)
		return helpers.ValidationErrorFromBind(err, input)
	}

	entity, err := h.service.Create(c.Context(), &input)
	if err != nil {
		slog.ErrorContext(c.Context(), "Error during create", "error", err)
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

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(httpx.NewResponse(map[string]bool{"deleted": true}))
}
