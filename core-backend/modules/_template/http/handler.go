package http

import (
	"net/http"

	"github.com/example/core-backend/modules/_template/app"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests.
type Handler struct {
	service *app.Service
}

// NewHandler constructs an HTTP handler.
func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes binds HTTP routes to chi.Router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/{id}", h.getByID)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	entity, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, httpx.Response{Data: entity})
}
