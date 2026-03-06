package domain

import "errors"

type ValidationError struct {
	Msg string
}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrNoSeats       = errors.New("no free seats")

	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingExpired      = errors.New("booking expired")
	ErrBookingInvalidState = errors.New("booking invalid state")
	ErrUserNotFound = errors.New("user not found")
	ErrAlreadyBooked = errors.New("user already has booking for this event")
	ErrConfirmationNotRequired = errors.New("confirmation is not required for this event")
)

func (e ValidationError) Error() string { return e.Msg }

func ErrValidation(msg string) error { return ValidationError{Msg: msg} }