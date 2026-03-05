package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
)

type bookSeatRequest struct {
	UserEmail string `json:"user_email"`
}

type confirmRequest struct {
	BookingID string `json:"booking_id"`
}

func (h *Handler) bookSeat(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req bookSeatRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	res, err := h.bookingSvc.BookSeat(r.Context(), eventID, req.UserEmail)
	if err != nil {
		if _, ok := err.(domain.ValidationError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, domain.ErrEventNotFound) {
			http.Error(w, "event not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrNoSeats) {
			http.Error(w, "no free seats", http.StatusConflict)
			return
		}

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) confirmBooking(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var req confirmRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	bookingID, err := uuid.Parse(req.BookingID)
	if err != nil {
		http.Error(w, "invalid booking_id", http.StatusBadRequest)
		return
	}

	err = h.bookingSvc.ConfirmBooking(r.Context(), eventID, bookingID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrBookingNotFound):
			http.Error(w, "booking not found", http.StatusNotFound)
			return
		case errors.Is(err, domain.ErrBookingExpired):
			http.Error(w, "booking expired", http.StatusConflict)
			return
		case errors.Is(err, domain.ErrBookingInvalidState):
			http.Error(w, "booking invalid state", http.StatusConflict)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "confirmed"})
}