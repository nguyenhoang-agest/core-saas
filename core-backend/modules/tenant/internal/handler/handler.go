package handler

import (
	"errors"
	"net/http"

	"github.com/example/core-backend/modules/tenant/internal/service"
	"github.com/example/core-backend/pkg/constants/response_code"
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
	var input service.CreateTenantRequest
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.ErrorCode(w, response_code.InvalidPayload)
		return
	}

	tenant, err := h.service.Create(r.Context(), &input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpx.Success(w, http.StatusCreated, tenant)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tenant, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpx.Success(w, http.StatusOK, tenant)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var respCode response_code.IResponseCode
	if errors.As(err, &respCode) {
		httpx.ErrorCode(w, respCode)
		return
	}

	httpx.ErrorCode(w, response_code.InternalServerGeneral)
}
