package wallet

import "errors"

// Predefined package errors
var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrNotFound         = errors.New("not found")
	ErrInvalidPIN       = errors.New("invalid pin code")
	ErrForbidden        = errors.New("forbidden")
)
