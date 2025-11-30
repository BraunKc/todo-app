package errors

import "errors"

var (
	ErrEmptyField   = errors.New("empty field")
	ErrTooLongField = errors.New("too long field")
	ErrInvalidField = errors.New("invalid field")
)
