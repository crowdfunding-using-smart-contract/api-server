package entity

import (
	"fund-o/api-server/pkg/pagination"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OldProject struct {
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

type Project struct {
	Base
	ProjectContractID string    `gorm:"type:varchar(255);not null"`
	Title             string    `gorm:"type:varchar(255);not null"`
	SubTitle          string    `gorm:"not null"`
	CategoryID        uuid.UUID `gorm:"not null"`
	Description       string
	Category          ProjectCategory `gorm:"foreignKey:CategoryID"`
	SubCategoryID     uuid.UUID
	SubCategory       ProjectSubCategory `gorm:"foreignKey:SubCategoryID"`
	Location          string             `gorm:"not null"`
	Image             string
	Ratings           []ProjectRating
	StartDate         time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	EndDate           time.Time `gorm:"not null"`
	OwnerID           uuid.UUID `gorm:"not null"`
	Owner             User      `gorm:"foreignKey:OwnerID"`
}

type ProjectDto struct {
	ID                string                 `json:"id"`
	ProjectContractID string                 `json:"project_contract_id"`
	Title             string                 `json:"title"`
	SubTitle          string                 `json:"sub_title"`
	Description       string                 `json:"description"`
	Category          *ProjectCategoryDto    `json:"category"`
	SubCategory       *ProjectSubCategoryDto `json:"sub_category"`
	Location          string                 `json:"location"`
	Image             string                 `json:"image"`
	Rating            float32                `json:"rating"`
	StartDate         string                 `json:"start_date"`
	EndDate           string                 `json:"end_date"`
	Owner             *UserDto               `json:"owner"`
	CreatedAt         string                 `json:"created_at"`
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
	ProjectContractID string                `form:"project_contract_id" binding:"required"`
	Title             string                `form:"title" binding:"required"`
	SubTitle          string                `form:"sub_title" binding:"required"`
	Description       string                `form:"description"`
	CategoryID        string                `form:"category_id" binding:"required"`
	SubCategoryID     string                `form:"sub_category_id" binding:"required"`
	Location          string                `form:"location" binding:"required"`
	Image             *multipart.FileHeader `form:"image" binding:"required"`
	EndDate           string                `form:"end_date" binding:"required"`
	OwnerID           string                `swaggerignore:"true"`
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
		ID:                p.ID.String(),
		ProjectContractID: p.ProjectContractID,
		Title:             p.Title,
		SubTitle:          p.SubTitle,
		Category:          p.Category.ToProjectCategoryDto(),
		SubCategory:       p.SubCategory.ToProjectSubCategoryDto(),
		Location:          p.Location,
		Rating:            rating,
		Image:             p.Image,
		Description:       p.Description,
		StartDate:         p.StartDate.Format(time.RFC3339),
		EndDate:           p.EndDate.Format(time.RFC3339),
		Owner:             p.Owner.ToUserDto(),
		CreatedAt:         p.CreatedAt.Format(time.RFC3339),
	}
}
