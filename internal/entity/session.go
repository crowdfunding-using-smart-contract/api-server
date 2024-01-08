package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID       uuid.UUID `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	UserAgent    string    `gorm:"not null"`
	ClientIP     string    `gorm:"not null"`
	IsBlocked    bool      `gorm:"default:false;not null"`
	ExpiredAt    time.Time `gorm:"not null"`
}

type SessionDto struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
}

// Secondary types

type SessionCreatePayload struct {
	ID           uuid.UUID
	UserID       string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiredAt    time.Time
} // @name SessionCreatePayload

// Parse functions

func (s *Session) ToSessionDto() *SessionDto {
	return &SessionDto{
		ID:           s.ID.String(),
		UserID:       s.UserID.String(),
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIP:     s.ClientIP,
		IsBlocked:    s.IsBlocked,
		ExpiredAt:    s.ExpiredAt,
	}
}
