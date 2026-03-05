package domain

import "errors"

type ValidationError struct {
	Msg string
}

var ErrEventNotFound = errors.New("event not found")

func (e ValidationError) Error() string { return e.Msg }

func ErrValidation(msg string) error { return ValidationError{Msg: msg} }