package apperrors

import "errors"

var (
	ErrInvalidMemberChannelLength = errors.New("invalid member channel length")
	ErrChannelNotFound            = errors.New("channel not found")
)
