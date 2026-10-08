package http

import (
	"net/http"

	"github.com/example/core-backend/modules/audit/app"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/logs", h.listLogs)
}

func (h *Handler) listLogs(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		tenantID = "default-tenant"
	}

	logs, err := h.service.ListByTenant(r.Context(), tenantID)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, httpx.Response{Data: logs})
}
