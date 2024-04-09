package entity

import (
	"fund-o/api-server/pkg/pagination"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Project struct {
	Base
	Title          string `gorm:"type:varchar(255);not null"`
	SubTitle       string `gorm:"not null"`
	Description    string
	CategoryID     uuid.UUID       `gorm:"not null"`
	Category       ProjectCategory `gorm:"foreignKey:CategoryID"`
	SubCategoryID  uuid.UUID
	SubCategory    ProjectSubCategory `gorm:"foreignKey:SubCategoryID"`
	Location       string             `gorm:"not null"`
	Image          string
	Ratings        []ProjectRating
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
	Description    string                 `json:"description"`
	Category       *ProjectCategoryDto    `json:"category"`
	SubCategory    *ProjectSubCategoryDto `json:"sub_category"`
	Location       string                 `json:"location"`
	Image          string                 `json:"image"`
	Rating         float32                `json:"rating"`
	TargetFunding  decimal.Decimal        `json:"target_funding"`
	CurrentFunding decimal.Decimal        `json:"current_funding"`
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

type ProjectListParams struct {
	pagination.PaginateOptions
	Query         string `form:"q"`
	CategoryID    string `form:"category"`
	SubCategoryID string `form:"sub_category"`
}

type ProjectCreatePayload struct {
	Title         string                `form:"title" binding:"required"`
	SubTitle      string                `form:"sub_title" binding:"required"`
	Description   string                `form:"description"`
	CategoryID    string                `form:"category_id" binding:"required"`
	SubCategoryID string                `form:"sub_category_id" binding:"required"`
	Location      string                `form:"location" binding:"required"`
	Image         *multipart.FileHeader `form:"image" binding:"required"`
	TargetFunding string                `form:"target_funding" binding:"required"`
	MonetaryUnit  string                `form:"monetary_unit"`
	EndDate       string                `form:"end_date" binding:"required"`
	LaunchDate    string                `form:"launch_date"`
	OwnerID       string                `swaggerignore:"true"`
}

type ProjectCreatePayloadSolidity struct {
	Title         string `form:"title" binding:"required"`
	SubTitle      string `form:"sub_title" binding:"required"`
	Description   string `form:"description"`
	CategoryID    string `form:"category_id" binding:"required"`
	SubCategoryID string `form:"sub_category_id" binding:"required"`
	Location      string `form:"location" binding:"required"`
	Image         string
	TargetFunding string `form:"target_funding" binding:"required"`
	MonetaryUnit  string `form:"monetary_unit"`
	EndDate       uint64 `form:"end_date" binding:"required"`
	LaunchDate    uint64 `form:"launch_date"`
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
		Location:       p.Location,
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
