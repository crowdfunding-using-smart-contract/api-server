package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID       uint      `gorm:"not null"`
	RefreshToken string    `gorm:"not null"`
	UserAgent    string    `gorm:"not null"`
	ClientIP     string    `gorm:"not null"`
	IsBlocked    bool      `gorm:"default:false;not null"`
	ExpiredAt    time.Time `gorm:"not null"`
}

// Secondary types

type SessionCreatePayload struct {
	ID           uuid.UUID
	UserID       uint
	RefreshToken string
	UserAgent    string
	ClientIP     string
	ExpiredAt    time.Time
} // @name SessionCreatePayload
