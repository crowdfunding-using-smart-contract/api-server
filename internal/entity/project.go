package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Project struct {
	Base
	Title          string          `gorm:"type:varchar(255);not null"`
	SubTitle       string          `gorm:"not null"`
	CategoryID     uuid.UUID       `gorm:"not null"`
	Category       ProjectCategory `gorm:"foreignKey:CategoryID"`
	SubCategoryID  uuid.UUID
	SubCategory    ProjectSubCategory `gorm:"foreignKey:SubCategoryID"`
	Image          string
	Description    string
	TargetFunding  decimal.Decimal `gorm:"type:decimal(32,16)"`
	CurrentFunding decimal.Decimal `gorm:"type:decimal(32,16)"`
	MonetaryUnit   string          `gorm:"default:'THB'"`
	StartDate      time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
	EndDate        time.Time       `gorm:"not null"`
	LaunchDate     time.Time
	OwnerID        uuid.UUID `gorm:"not null"`
	Owner          User      `gorm:"foreignKey:OwnerID"`
}

type ProjectDto struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	SubTitle       string                 `json:"sub_title"`
	Category       *ProjectCategoryDto    `json:"category"`
	SubCategory    *ProjectSubCategoryDto `json:"sub_category"`
	Image          string                 `json:"image"`
	Description    string                 `json:"description"`
	TargetFunding  decimal.Decimal        `json:"target_amount"`
	CurrentFunding decimal.Decimal        `json:"current_amount"`
	MonetaryUnit   string                 `json:"monetary_unit"`
	StartDate      string                 `json:"start_date"`
	EndDate        string                 `json:"end_date"`
	LaunchDate     string                 `json:"launch_date"`
	Owner          *UserDto               `json:"owner"`
	CreatedAt      string                 `json:"created_at"`
} // @name Project

// Secondary types

type ProjectCreatePayload struct {
	Title         string          `json:"title" binding:"required"`
	SubTitle      string          `json:"sub_title" binding:"required"`
	CategoryID    string          `json:"category_id" binding:"required"`
	SubCategoryID string          `json:"sub_category_id" binding:"required"`
	Image         string          `json:"image"`
	Description   string          `json:"description"`
	TargetFunding decimal.Decimal `json:"target_amount" binding:"required"`
	MonetaryUnit  string          `json:"monetary_unit"`
	EndDate       string          `json:"end_date" binding:"required"`
	LaunchDate    string          `json:"launch_date"`
	OwnerID       string
}

// Parse functions

func (p *Project) ToProjectDto() *ProjectDto {
	return &ProjectDto{
		ID:             p.ID.String(),
		Title:          p.Title,
		SubTitle:       p.SubTitle,
		Category:       p.Category.ToProjectCategoryDto(),
		SubCategory:    p.SubCategory.ToProjectSubCategoryDto(),
		Image:          p.Image,
		Description:    p.Description,
		TargetFunding:  p.TargetFunding,
		CurrentFunding: p.CurrentFunding,
		MonetaryUnit:   p.MonetaryUnit,
		StartDate:      p.StartDate.Format(time.RFC3339),
		EndDate:        p.EndDate.Format(time.RFC3339),
		LaunchDate:     p.LaunchDate.Format(time.RFC3339),
		Owner:          p.Owner.ToUserDto(),
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
	}
}
