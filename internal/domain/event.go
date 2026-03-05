package domain

import "time"

type Event struct {
	ID                string
	Title             string
	StartsAt          time.Time
	Capacity          int
	RequiresPayment   bool
	BookingTTLSeconds int
	CreatedAt         time.Time
}