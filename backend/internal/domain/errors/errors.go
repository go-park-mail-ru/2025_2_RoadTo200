package errors

import "errors"

var (
	// User errors
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("invalid email")

	// Password errors
	ErrPasswordTooShort   = errors.New("password is too short")
	ErrPasswordsDontMatch = errors.New("passwords don't match")
	ErrInvalidPassword    = errors.New("invalid password")

	// Session errors
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")

	// Auth errors
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")

	// General errors
	ErrInternalError = errors.New("internal error")
)
