package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/MyNameIsWhaaat/event-booker/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
		INSERT INTO bookings (event_id, user_id, user_email, status, expires_at)
		VALUES ($1, $2, $3, 'pending', $4)
		RETURNING id
	`

	var id uuid.UUID
	err := tx.QueryRowContext(ctx, q,
		b.EventID,
		b.UserID,
		b.UserEmail,
		b.ExpiresAt,
	).Scan(&id)

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "bookings_event_user_active_uidx" {
				return uuid.Nil, domain.ErrAlreadyBooked
			}
		}

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

func (r *BookingRepository) CancelExpired(ctx context.Context, now time.Time) ([]domain.Booking, error) {
	const q = `
		UPDATE bookings
		SET status = 'cancelled',
		    cancelled_at = $1
		WHERE status = 'pending'
		  AND expires_at <= $1
		RETURNING
			id,
			event_id,
			user_id,
			user_email,
			status,
			created_at,
			expires_at,
			confirmed_at,
			cancelled_at
	`

	rows, err := r.db.QueryContext(ctx, q, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Booking

	for rows.Next() {
		var b domain.Booking

		err := rows.Scan(
			&b.ID,
			&b.EventID,
			&b.UserID,
			&b.UserEmail,
			&b.Status,
			&b.CreatedAt,
			&b.ExpiresAt,
			&b.ConfirmedAt,
			&b.CancelledAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *BookingRepository) GetEventStats(ctx context.Context, eventID uuid.UUID) (repository.EventBookingStats, error) {
	const q = `
		SELECT
			COUNT(*) FILTER (WHERE status = 'pending')   AS pending,
			COUNT(*) FILTER (WHERE status = 'confirmed') AS confirmed
		FROM bookings
		WHERE event_id = $1
	`

	var s repository.EventBookingStats
	if err := r.db.QueryRowContext(ctx, q, eventID).Scan(&s.Pending, &s.Confirmed); err != nil {
		return repository.EventBookingStats{}, err
	}
	return s, nil
}

func (r *BookingRepository) CreateConfirmed(ctx context.Context, tx *sql.Tx, b domain.Booking, now time.Time) (uuid.UUID, error) {
	const q = `
		INSERT INTO bookings (event_id, user_id, user_email, status, expires_at, confirmed_at)
		VALUES ($1, $2, $3, 'confirmed', $4, $4)
		RETURNING id
	`

	var id uuid.UUID
	if err := tx.QueryRowContext(ctx, q, b.EventID, b.UserID, b.UserEmail, now).Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *BookingRepository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Booking, error) {
	const q = `
		SELECT
			id,
			event_id,
			user_id,
			user_email,
			status,
			created_at,
			expires_at,
			confirmed_at,
			cancelled_at
		FROM bookings
		WHERE event_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Booking

	for rows.Next() {
		var b domain.Booking

		err := rows.Scan(
			&b.ID,
			&b.EventID,
			&b.UserID,
			&b.UserEmail,
			&b.Status,
			&b.CreatedAt,
			&b.ExpiresAt,
			&b.ConfirmedAt,
			&b.CancelledAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}