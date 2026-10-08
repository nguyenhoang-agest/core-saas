package http

import (
	"net/http"

	"github.com/example/core-backend/modules/example/app"
	"github.com/example/core-backend/pkg/httpx"
	"github.com/example/core-backend/pkg/tenantctx"
	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for example items.
type Handler struct {
	service *app.Service
}

// NewHandler constructs an HTTP handler for the example module.
func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes attaches module endpoints to a chi.Router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Use(h.extractTenant)

	r.Post("/items", h.createItem)
	r.Get("/items", h.listItems)
	r.Get("/items/{id}", h.getItem)
}

func (h *Handler) extractTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "default-tenant"
		}
		ctx := tenantctx.WithContext(r.Context(), tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {
	var input app.CreateItemInput
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, err)
		return
	}

	item, err := h.service.Create(r.Context(), input)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, httpx.Response{Data: item})
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, httpx.Response{Data: item})
}

func (h *Handler) listItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, httpx.Response{Data: items})
}
