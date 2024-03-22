package apperrors

import "errors"

var (
	ErrAlreadyRatedProject = errors.New("user already rated this project")
)
