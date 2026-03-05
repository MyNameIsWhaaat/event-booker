package repository

import (
	"context"
	// "database/sql"
	// "time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/google/uuid"
)

type EventRepository interface {
	Create(ctx context.Context, e domain.Event) (uuid.UUID, error)
	// GetByID(ctx context.Context, id uuid.UUID) (domain.Event, error)
	// GetByIDForUpdate(ctx context.Context, q Querier, id uuid.UUID) (domain.Event, error)
}

type BookingRepository interface {
	// CountActiveByEvent(ctx context.Context, q Querier, eventID uuid.UUID) (int, error)
	// CreatePending(ctx context.Context, q Querier, b domain.Booking) (uuid.UUID, error)
	// ConfirmPending(ctx context.Context, q Querier, eventID, bookingID uuid.UUID, now time.Time) (bool, error)
	// CancelExpired(ctx context.Context, now time.Time) (int, error)
}

type Querier interface {
	// ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Transactor interface {
	// WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}

type Repositories struct {
	Events   EventRepository
	Bookings BookingRepository
}