package http

import (
	"net/http"

	"github.com/MyNameIsWhaaat/event-booker/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	eventSvc   service.EventService
	bookingSvc service.BookingService
}

func New(eventSvc service.EventService, bookingSvc service.BookingService) *Handler {
	return &Handler{eventSvc: eventSvc, bookingSvc: bookingSvc}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Post("/events", h.createEvent)
	r.Get("/events/{id}", h.getEvent)

	return r
}
