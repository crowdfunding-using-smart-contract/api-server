package entity

import (
	"time"
)

type VerifyEmail struct {
	Base
	Email      string    `gorm:"not null"`
	SecretCode string    `gorm:"not null"`
	IsUsed     bool      `gorm:"not null;default:false"`
	ExpiredAt  time.Time `gorm:"not null;default:now() + interval '15 minutes'"`
}

type VerifyEmailDto struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
	IsUsed     bool   `json:"is_used"`
	ExpiredAt  string `json:"expired_at"`
} // @name VerifyEmail

// Secondary types

type VerifyEmailCreatePayload struct {
	Email      string `json:"email" binding:"required"`
	SecretCode string `json:"secret_code" binding:"required"`
} // @name VerifyEmailCreatePayload

// Parse functions

func (v *VerifyEmail) ToVerifyEmailDto() *VerifyEmailDto {
	return &VerifyEmailDto{
		ID:         v.ID.String(),
		Email:      v.Email,
		SecretCode: v.SecretCode,
		IsUsed:     v.IsUsed,
		ExpiredAt:  v.ExpiredAt.Format(time.RFC3339),
	}
}
