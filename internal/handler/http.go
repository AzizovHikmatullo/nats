package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"nats/internal/model"
	"nats/internal/nats"
)

type Handler struct {
	js *nats.Client
}

func New(js *nats.Client) *Handler {
	return &Handler{js: js}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/coordinates", h.handleCoordinates)
}

func (h *Handler) handleCoordinates(w http.ResponseWriter, r *http.Request) {
	var req model.LocationEvent
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.js.Publish(ctx, req); err != nil {
		http.Error(w, "failed to publish", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
