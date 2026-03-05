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
	bookings repository.BookingRepository
}

type CreateEventRequest struct {
	Title             string
	StartsAt          time.Time
	Capacity          int
	RequiresPayment   bool
	BookingTTLSeconds int
}

type EventDetails struct {
	Event domain.Event
	Stats struct {
		Pending   int `json:"pending"`
		Confirmed int `json:"confirmed"`
		FreeSeats int `json:"free_seats"`
	} `json:"stats"`
}

func NewEventService(events repository.EventRepository, bookings repository.BookingRepository) EventService {
	return &eventService{events: events, bookings: bookings}
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

// func (s *eventService) GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error) {
// 	return s.events.GetByID(ctx, id)
// }

func (s *eventService) GetEventDetails(ctx context.Context, id uuid.UUID) (EventDetails, error) {
	ev, err := s.events.GetByID(ctx, id)
	if err != nil {
		return EventDetails{}, err
	}

	stats, err := s.bookings.GetEventStats(ctx, id)
	if err != nil {
		return EventDetails{}, err
	}

	free := ev.Capacity - stats.Pending - stats.Confirmed
	if free < 0 {
		free = 0
	}

	var d EventDetails
	d.Event = ev
	d.Stats.Pending = stats.Pending
	d.Stats.Confirmed = stats.Confirmed
	d.Stats.FreeSeats = free
	return d, nil
}