package service

import (
	"context"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(ctx context.Context, req CreateEventRequest) (uuid.UUID, error)
}

type BookingService interface {
	
}

type Services struct {
	Events   EventService
	Bookings BookingService
}