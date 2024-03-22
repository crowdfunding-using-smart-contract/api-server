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
	Ratings        []ProjectRating
	Image          string
	Description    string
	TargetFunding  decimal.Decimal `gorm:"type:decimal(32,16)"`
	CurrentFunding decimal.Decimal `gorm:"type:decimal(32,16)"`
	MonetaryUnit   string          `gorm:"default:'THB'"`
	StartDate      time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
	EndDate        time.Time       `gorm:"not null"`
	LaunchDate     time.Time       `gorm:"default:null"`
	OwnerID        uuid.UUID       `gorm:"not null"`
	Owner          User            `gorm:"foreignKey:OwnerID"`
}

type ProjectDto struct {
	ID             string                 `json:"id"`
	Title          string                 `json:"title"`
	SubTitle       string                 `json:"sub_title"`
	Category       *ProjectCategoryDto    `json:"category"`
	SubCategory    *ProjectSubCategoryDto `json:"sub_category"`
	Rating         float32                `json:"rating"`
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

type ProjectRating struct {
	Base
	Rating    float32 `gorm:"not null"`
	ProjectID uuid.UUID
	UserID    uuid.UUID
} // @name ProjectRating

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
	OwnerID       string          `swaggerignore:"true"`
}

type ProjectRatingCreatePayload struct {
	Rating    float32 `json:"rating" binding:"required,gte=0,lte=5"`
	ProjectID string  `json:"project_id" binding:"required" swaggerignore:"true"`
	UserID    string  `swaggerignore:"true"`
}

// Parse functions

func (p *Project) ToProjectDto() *ProjectDto {
	rating := float32(0)

	if len(p.Ratings) > 0 {
		totalRating := float32(0)
		for _, r := range p.Ratings {
			totalRating += r.Rating
		}
		rating = totalRating / float32(len(p.Ratings))
	}

	return &ProjectDto{
		ID:             p.ID.String(),
		Title:          p.Title,
		SubTitle:       p.SubTitle,
		Category:       p.Category.ToProjectCategoryDto(),
		SubCategory:    p.SubCategory.ToProjectSubCategoryDto(),
		Rating:         rating,
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
