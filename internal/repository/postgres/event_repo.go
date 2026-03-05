package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Event, error) {
	const q = `
		SELECT id, title, starts_at, capacity, requires_payment, booking_ttl_seconds, created_at
		FROM events
		WHERE id = $1
	`

	var e domain.Event
	var eid uuid.UUID
	var createdAt time.Time

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&eid,
		&e.Title,
		&e.StartsAt,
		&e.Capacity,
		&e.RequiresPayment,
		&e.BookingTTLSeconds,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, err
	}

	e.ID = eid.String()
	e.CreatedAt = createdAt
	return e, nil
}