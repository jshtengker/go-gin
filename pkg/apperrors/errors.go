package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrConflict     = errors.New("record already exists")
	ErrInvalidInput = errors.New("invalid input data")
)
