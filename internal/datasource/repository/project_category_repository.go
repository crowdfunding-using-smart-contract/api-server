package repository

import (
	"fund-o/api-server/internal/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProjectCategoryRepository interface {
	FindAll() ([]entity.ProjectCategory, error)
}

type projectCategoryRepository struct {
	db     *gorm.DB
	logger *log.Entry
}

func NewProjectCategoryRepository(db *gorm.DB) ProjectCategoryRepository {
	logger := log.WithFields(log.Fields{
		"module": "project_category_repository",
	})
	return &projectCategoryRepository{db, logger}
}

func (repo *projectCategoryRepository) FindAll() ([]entity.ProjectCategory, error) {
	var categories []entity.ProjectCategory
	if result := repo.db.Preload("SubCategories").Find(&categories); result.Error != nil {
		repo.logger.Errorf("Failed to list project categories: %v", result.Error)
		return nil, result.Error
	}

	return categories, nil
}
