package apperrors

import "errors"

var (
	ErrPasswordAndConfirmationNotMatch = errors.New("password and password confirmation does not match")
	ErrHashPassword                    = errors.New("failed to hash password")
	ErrInvalidBirthDateFormat          = errors.New("invalid birth date format")
)
