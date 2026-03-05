package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/MyNameIsWhaaat/event-booker/internal/service"
)

type createEventRequest struct {
	Title             string    `json:"title"`
	StartsAt          time.Time `json:"starts_at"`
	Capacity          int       `json:"capacity"`
	RequiresPayment   bool      `json:"requires_payment"`
	BookingTTLSeconds int       `json:"booking_ttl_seconds"`
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	id, err := h.eventSvc.CreateEvent(r.Context(), service.CreateEventRequest{
		Title:             req.Title,
		StartsAt:          req.StartsAt,
		Capacity:          req.Capacity,
		RequiresPayment:   req.RequiresPayment,
		BookingTTLSeconds: req.BookingTTLSeconds,
	})
	if err != nil {
		if _, ok := err.(domain.ValidationError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
}