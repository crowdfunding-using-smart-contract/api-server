package entity

import "github.com/google/uuid"

type Category int

type ProjectCategory struct {
	Base
	Name          string               `gorm:"type:varchar(255);not null"`
	SubCategories []ProjectSubCategory `gorm:"foreignKey:CategoryID"`
}

type ProjectSubCategory struct {
	Base
	Name       string `gorm:"type:varchar(255);not null"`
	CategoryID uuid.UUID
}

type ProjectCategoryDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name ProjectCategory

type ProjectSubCategoryDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name ProjectSubCategoryfc3c774c-ea74-4886-a70c-05ff62cc62c0

// Parse functions

func (p *ProjectCategory) ToProjectCategoryDto() *ProjectCategoryDto {
	return &ProjectCategoryDto{
		ID:   p.ID.String(),
		Name: p.Name,
	}
}

func (p *ProjectSubCategory) ToProjectSubCategoryDto() *ProjectSubCategoryDto {
	return &ProjectSubCategoryDto{
		ID:   p.ID.String(),
		Name: p.Name,
	}
}
