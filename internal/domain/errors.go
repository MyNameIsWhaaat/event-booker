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
)

func (e ValidationError) Error() string { return e.Msg }

func ErrValidation(msg string) error { return ValidationError{Msg: msg} }