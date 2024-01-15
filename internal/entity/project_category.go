package entity

import "github.com/google/uuid"

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
	ID            string `json:"id"`
	Name          string `json:"name"`
	SubCategories []ProjectSubCategoryDto
}

type ProjectSubCategoryDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
