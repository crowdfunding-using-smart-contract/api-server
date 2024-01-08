package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProjectStatus int

const (
	Open ProjectStatus = iota + 1
	Approved
	Reverted
	Deleted
	PaidOut
)

type Project struct {
	Base
	Title         string `gorm:"type:varchar(255);not null"`
	Description   string `gorm:"not null"`
	Image         string
	TargetAmount  decimal.Decimal `gorm:"type:decimal(32,16)"`
	CurrentAmount decimal.Decimal `gorm:"type:decimal(32,16)"`
	OwnerID       uuid.UUID       `gorm:"not null"`
	Owner         User            `gorm:"foreignKey:OwnerID"`
}

type ProjectDto struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Image         string          `json:"image"`
	TargetAmount  decimal.Decimal `json:"target_amount"`
	CurrentAmount decimal.Decimal `json:"current_amount"`
	Owner         *UserDto        `json:"owner"`
	CreatedAt     string          `json:"created_at"`
} // @name Project

// Secondary types

type ProjectCreatePayload struct {
	Title        string          `json:"title" binding:"required"`
	Description  string          `json:"description" binding:"required"`
	Image        string          `json:"image"`
	TargetAmount decimal.Decimal `json:"target_amount" binding:"required"`
	OwnerID      string
}

// Parse functions

func (p *Project) ToProjectDto() *ProjectDto {
	return &ProjectDto{
		ID:            p.ID.String(),
		Title:         p.Title,
		Description:   p.Description,
		Image:         p.Image,
		TargetAmount:  p.TargetAmount,
		CurrentAmount: p.CurrentAmount,
		Owner:         p.Owner.ToUserDto(),
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
	}
}
