package apperrors

import "errors"

var (
	ErrNotFound              = errors.New("record not found")
	ErrConflict              = errors.New("record already exists")
	ErrInvalidInput          = errors.New("invalid input data")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInvalidCredentials    = errors.New("invalid email/username or password")
	ErrEmailAlreadyExists    = errors.New("email is already registered")
	ErrUsernameAlreadyExists = errors.New("username is already taken")
	ErrInvalidToken          = errors.New("token is invalid or expired")
	ErrTokenExpired          = errors.New("token has expired")
	ErrTokenRevoked          = errors.New("token has been revoked")
	ErrUserInactive          = errors.New("user account is inactive")
)
