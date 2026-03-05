package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/MyNameIsWhaaat/event-booker/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	details, err := h.eventSvc.GetEventDetails(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(details)
}

func (h *Handler) listEvents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}

	items, err := h.eventSvc.ListEvents(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"items": items,
		"limit": limit,
		"offset": offset,
	})
}