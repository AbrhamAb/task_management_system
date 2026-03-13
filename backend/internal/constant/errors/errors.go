package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTaskNotFound       = errors.New("task not found")
)
