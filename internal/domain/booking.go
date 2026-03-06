package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingConfirmed BookingStatus = "confirmed"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID          string
	EventID     string
	UserID      uuid.UUID
	UserEmail   string
	Status      BookingStatus
	CreatedAt   time.Time
	ExpiresAt   time.Time
	ConfirmedAt *time.Time
	CancelledAt *time.Time
}
