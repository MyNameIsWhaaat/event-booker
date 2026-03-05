package postgres

import (
	"context"
	"database/sql"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, e domain.Event) (uuid.UUID, error) {
	const q = `
		INSERT INTO events (title, starts_at, capacity, requires_payment, booking_ttl_seconds)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(
		ctx, q,
		e.Title, e.StartsAt, e.Capacity, e.RequiresPayment, e.BookingTTLSeconds,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}