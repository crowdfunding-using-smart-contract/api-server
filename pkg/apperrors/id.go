package apperrors

import "errors"

var (
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidProjectID = errors.New("invalid project id")
)
