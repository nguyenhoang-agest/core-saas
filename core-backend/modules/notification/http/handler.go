package http

import (
	"net/http"

	"github.com/example/core-backend/modules/notification/app"
	"github.com/example/core-backend/modules/notification/domain"
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
	r.Post("/send", h.send)
}

func (h *Handler) send(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Recipient string `json:"recipient"`
		Subject   string `json:"subject"`
		Body      string `json:"body"`
		Channel   string `json:"channel"`
	}
	var input req
	if err := httpx.DecodeJSON(r, &input); err != nil {
		httpx.Error(w, err)
		return
	}

	msg := &domain.Message{
		Recipient: input.Recipient,
		Subject:   input.Subject,
		Body:      input.Body,
		Channel:   domain.Channel(input.Channel),
	}

	if err := h.service.Dispatch(r.Context(), msg); err != nil {
		httpx.Error(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]string{"status": "queued"})
}
