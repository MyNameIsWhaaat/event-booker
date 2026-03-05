package domain

import "errors"

var (
	ErrEventNotFound   = errors.New("event not found")
	ErrNoSeats         = errors.New("no free seats")
	ErrBookingNotFound = errors.New("booking not found")
	ErrBookingExpired  = errors.New("booking expired")
	ErrBookingInvalid  = errors.New("booking invalid state")
)