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
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	SubCategories []ProjectSubCategoryDto `json:"sub_categories"`
} // @name ProjectCategory

type ProjectSubCategoryDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name ProjectSubCategory

// Parse functions

func (p *ProjectCategory) ToProjectCategoryDto() *ProjectCategoryDto {
	subCategories := make([]ProjectSubCategoryDto, len(p.SubCategories))
	for i, subCategory := range p.SubCategories {
		subCategories[i] = *subCategory.ToProjectSubCategoryDto()
	}

	return &ProjectCategoryDto{
		ID:            p.ID.String(),
		Name:          p.Name,
		SubCategories: subCategories,
	}
}

func (p *ProjectSubCategory) ToProjectSubCategoryDto() *ProjectSubCategoryDto {
	return &ProjectSubCategoryDto{
		ID:   p.ID.String(),
		Name: p.Name,
	}
}
