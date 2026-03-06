package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/MyNameIsWhaaat/event-booker/internal/domain"
	"github.com/MyNameIsWhaaat/event-booker/internal/repository"
	"github.com/google/uuid"
)

type bookingService struct {
	tx       repository.Transactor
	events   repository.EventRepository
	bookings repository.BookingRepository
	users    repository.UserRepository
}

type BookSeatResult struct {
	BookingID uuid.UUID `json:"booking_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewBookingService(tx repository.Transactor, events repository.EventRepository, bookings repository.BookingRepository, users repository.UserRepository) BookingService {
	return &bookingService{
		tx:       tx,
		events:   events,
		bookings: bookings,
		users:    users,
	}
}

func (s *bookingService) BookSeat(ctx context.Context, eventID uuid.UUID, userEmail string) (BookSeatResult, error) {
	if userEmail == "" {
		return BookSeatResult{}, domain.ErrValidation("user_email is required")
	}

	now := time.Now().UTC()

	var res BookSeatResult
	err := s.tx.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		ev, err := s.events.GetByIDForUpdate(ctx, tx, eventID)
		if err != nil {
			return err
		}

		active, err := s.bookings.CountActiveByEvent(ctx, tx, eventID)
		if err != nil {
			return err
		}
		if active >= ev.Capacity {
			return domain.ErrNoSeats
		}

		user, err := s.users.GetOrCreateByEmail(ctx, userEmail)
		if err != nil {
			return err
		}

		if !ev.RequiresPayment {
			bookingID, err := s.bookings.CreateConfirmed(ctx, tx, domain.Booking{
				EventID:   eventID.String(),
				UserID:    user.ID,
				UserEmail: user.Email,
			}, now)
			if err != nil {
				return err
			}

			res = BookSeatResult{
				BookingID: bookingID,
				Status:    "confirmed",
				ExpiresAt: now,
			}
			return nil
		}

		expiresAt := now.Add(time.Duration(ev.BookingTTLSeconds) * time.Second)
		bookingID, err := s.bookings.CreatePending(ctx, tx, domain.Booking{
			EventID:   eventID.String(),
			UserID:    user.ID,
			UserEmail: user.Email,
			ExpiresAt: expiresAt,
		})
		if err != nil {
			return err
		}

		res = BookSeatResult{
			BookingID: bookingID,
			Status:    "pending",
			ExpiresAt: expiresAt,
		}
		return nil
	})
	if err != nil {
		return BookSeatResult{}, err
	}

	return res, nil
}

func (s *bookingService) ConfirmBooking(ctx context.Context, eventID, bookingID uuid.UUID) error {
	ev, err := s.events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if !ev.RequiresPayment {
		return domain.ErrConfirmationNotRequired
	}

	now := time.Now().UTC()

	return s.tx.WithinTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return s.bookings.ConfirmPending(ctx, tx, eventID, bookingID, now)
	})
}

func (s *bookingService) CancelExpired(ctx context.Context) ([]domain.Booking, error) {
	now := time.Now().UTC()
	return s.bookings.CancelExpired(ctx, now)
}

func (s *bookingService) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Booking, error) {
	return s.bookings.ListByEvent(ctx, eventID)
}