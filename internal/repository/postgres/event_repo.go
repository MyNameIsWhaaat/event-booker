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

func (r *EventRepository) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id uuid.UUID) (domain.Event, error) {
	const q = `
		SELECT id, title, starts_at, capacity, requires_payment, booking_ttl_seconds, created_at
		FROM events
		WHERE id = $1
		FOR UPDATE
	`

	var e domain.Event
	var eid uuid.UUID
	var createdAt time.Time

	err := tx.QueryRowContext(ctx, q, id).Scan(
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

func (r *EventRepository) List(ctx context.Context, limit, offset int) ([]domain.Event, error) {
	const q = `
		SELECT id, title, starts_at, capacity, requires_payment, booking_ttl_seconds, created_at
		FROM events
		ORDER BY starts_at ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Event, 0, limit)

	for rows.Next() {
		var (
			e        domain.Event
			id       uuid.UUID
			createdAt time.Time
		)

		if err := rows.Scan(
			&id,
			&e.Title,
			&e.StartsAt,
			&e.Capacity,
			&e.RequiresPayment,
			&e.BookingTTLSeconds,
			&createdAt,
		); err != nil {
			return nil, err
		}

		e.ID = id.String()
		e.CreatedAt = createdAt
		out = append(out, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}