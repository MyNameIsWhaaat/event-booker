package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/google/uuid"
)

type EventBookingStats struct {
	Pending   int
	Confirmed int
}

type EventRepository interface {
	Create(ctx context.Context, e domain.Event) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Event, error)
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id uuid.UUID) (domain.Event, error)
	List(ctx context.Context, limit, offset int) ([]domain.Event, error)
}

type BookingRepository interface {
	CountActiveByEvent(ctx context.Context, tx *sql.Tx, eventID uuid.UUID) (int, error)
	CreatePending(ctx context.Context, tx *sql.Tx, b domain.Booking) (uuid.UUID, error)
	ConfirmPending(ctx context.Context, tx *sql.Tx, eventID, bookingID uuid.UUID, now time.Time) error
	CancelExpired(ctx context.Context, now time.Time) (int, error)
	GetEventStats(ctx context.Context, eventID uuid.UUID) (EventBookingStats, error)
	CreateConfirmed(ctx context.Context, tx *sql.Tx, b domain.Booking, now time.Time) (uuid.UUID, error)
}

type Repositories struct {
	Events   EventRepository
	Bookings BookingRepository
}