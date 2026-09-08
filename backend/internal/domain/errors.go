package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrDuplicateEmail     = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
