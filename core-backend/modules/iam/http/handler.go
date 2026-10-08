package http

import (
	"net/http"

	"github.com/example/core-backend/modules/iam/app"
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
	r.Post("/register", h.register)
	r.Post("/login", h.login)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input req
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, err)
		return
	}

	user, err := h.service.Register(r.Context(), input.Email, input.Password)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, httpx.Response{Data: user})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input req
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, err)
		return
	}

	token, err := h.service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]string{"token": token})
}
