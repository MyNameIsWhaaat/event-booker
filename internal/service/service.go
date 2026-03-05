package service

import (
	"context"

	// "github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, req CreateEventRequest) (uuid.UUID, error)
	// GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error)
	GetEventDetails(ctx context.Context, id uuid.UUID) (EventDetails, error)
}

type BookingService interface {
	BookSeat(ctx context.Context, eventID uuid.UUID, userEmail string) (BookSeatResult, error)
	ConfirmBooking(ctx context.Context, eventID, bookingID uuid.UUID) error
	CancelExpired(ctx context.Context) (int, error)
}

type Services struct {
	Events   EventService
	Bookings BookingService
}