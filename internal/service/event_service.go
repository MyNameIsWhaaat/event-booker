package service

import (
	"context"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/MyNameIsWhaaat/event-booker/internal/repository"
	"github.com/google/uuid"
)

type eventService struct {
	events repository.EventRepository
}

type CreateEventRequest struct {
	Title             string
	StartsAt          time.Time
	Capacity          int
	RequiresPayment   bool
	BookingTTLSeconds int
}

func NewEventService(events repository.EventRepository) EventService {
	return &eventService{events: events}
}

func (s *eventService) CreateEvent(ctx context.Context, req CreateEventRequest) (uuid.UUID, error) {
	if req.Title == "" {
		return uuid.Nil, domain.ErrValidation("title is required")
	}
	if req.Capacity <= 0 {
		return uuid.Nil, domain.ErrValidation("capacity must be > 0")
	}
	if req.BookingTTLSeconds <= 0 {
		return uuid.Nil, domain.ErrValidation("booking_ttl_seconds must be > 0")
	}

	e := domain.Event{
		Title:             req.Title,
		StartsAt:          req.StartsAt,
		Capacity:          req.Capacity,
		RequiresPayment:   req.RequiresPayment,
		BookingTTLSeconds: req.BookingTTLSeconds,
	}

	return s.events.Create(ctx, e)
}