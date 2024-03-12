package repository

import (
	"fund-o/api-server/internal/entity"
	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ProjectCategoryRepository interface {
	FindAll() ([]entity.ProjectCategory, error)
}

type projectCategoryRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewProjectCategoryRepository(db *gorm.DB) ProjectCategoryRepository {
	logger := log.With().Str("module", "project_category_repository").Logger()
	return &projectCategoryRepository{db, logger}
}

func (repo *projectCategoryRepository) FindAll() ([]entity.ProjectCategory, error) {
	var categories []entity.ProjectCategory
	if result := repo.db.Preload("SubCategories").Find(&categories); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list project categories")
		return nil, result.Error
	}

	return categories, nil
}
