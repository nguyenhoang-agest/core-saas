package http

import (
	"errors"
	"net/http"

	"github.com/example/core-backend/modules/tenant/domain/dto"
	"github.com/example/core-backend/modules/tenant/domain/exception"
	"github.com/example/core-backend/modules/tenant/service"
	pkgerrors "github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service.ITenantService
}

func NewHandler(service service.ITenantService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.create)
	r.Get("/{id}", h.getByID)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateTenantRequest
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, err)
		return
	}

	tenant, err := h.service.Create(r.Context(), &input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, httpx.Response{Data: tenant})
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, httpx.Response{Data: tenant})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, exception.ErrNotFound):
		httpx.Error(w, pkgerrors.NotFound("tenant not found", err))
	case errors.Is(err, exception.ErrInvalidInput):
		httpx.Error(w, pkgerrors.InvalidArgument(err.Error(), err))
	case errors.Is(err, exception.ErrAlreadyExists):
		httpx.Error(w, pkgerrors.Conflict("tenant already exists", err))
	default:
		httpx.Error(w, pkgerrors.Internal("internal server error", err))
	}
}
