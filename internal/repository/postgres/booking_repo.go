package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/google/uuid"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CountActiveByEvent(ctx context.Context, tx *sql.Tx, eventID uuid.UUID) (int, error) {
	const q = `
		SELECT COUNT(*)
		FROM bookings
		WHERE event_id = $1
		  AND status IN ('pending', 'confirmed')
	`

	var n int
	if err := tx.QueryRowContext(ctx, q, eventID).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (r *BookingRepository) CreatePending(ctx context.Context, tx *sql.Tx, b domain.Booking) (uuid.UUID, error) {
	const q = `
		INSERT INTO bookings (event_id, user_email, status, expires_at)
		VALUES ($1, $2, 'pending', $3)
		RETURNING id
	`

	var id uuid.UUID
	if err := tx.QueryRowContext(ctx, q, b.EventID, b.UserEmail, b.ExpiresAt).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *BookingRepository) ConfirmPending(
	ctx context.Context,
	tx *sql.Tx,
	eventID, bookingID uuid.UUID,
	now time.Time,
) error {
	const upd = `
		UPDATE bookings
		SET status = 'confirmed', confirmed_at = $3
		WHERE id = $1
		  AND event_id = $2
		  AND status = 'pending'
		  AND expires_at > $3
		RETURNING id
	`

	var id uuid.UUID
	err := tx.QueryRowContext(ctx, upd, bookingID, eventID, now).Scan(&id)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	const sel = `
		SELECT status, expires_at
		FROM bookings
		WHERE id = $1 AND event_id = $2
	`

	var status string
	var expiresAt time.Time
	err = tx.QueryRowContext(ctx, sel, bookingID, eventID).Scan(&status, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrBookingNotFound
	}
	if err != nil {
		return err
	}

	if !expiresAt.After(now) {
		return domain.ErrBookingExpired
	}
	return domain.ErrBookingInvalidState
}