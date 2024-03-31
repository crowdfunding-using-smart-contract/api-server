package repository

import (
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/pagination"
	"github.com/rs/zerolog"

	"github.com/google/uuid"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	FindAll(findOptions pagination.PaginateFindOptions) []entity.Project
	Count() int64
	Create(project *entity.Project) (*entity.Project, error)
	FindByID(projectID uuid.UUID) (*entity.Project, error)
	FindAllByOwnerID(ownerID uuid.UUID) ([]entity.Project, error)
	CreateProjectRating(rating *entity.ProjectRating) (*entity.ProjectRating, error)
	FindProjectRating(userID uuid.UUID, projectID uuid.UUID) (*entity.ProjectRating, error)
}

type projectRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	logger := log.With().Str("module", "project_repository").Logger()
	return &projectRepository{db, logger}
}

func (repo *projectRepository) FindAll(findOptions pagination.PaginateFindOptions) (projects []entity.Project) {
	result := repo.db.
		Limit(findOptions.Limit).
		Offset(findOptions.Skip).
		Preload("Category").
		Preload("SubCategory").
		Preload("Owner").
		Preload("Ratings").
		Find(&projects)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list projects")
		return
	}

	return projects
}

func (repo *projectRepository) Count() int64 {
	var count int64
	if result := repo.db.Model(&entity.Project{}).Count(&count); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to count projects")
		return 0
	}

	return count
}

func (repo *projectRepository) Create(project *entity.Project) (*entity.Project, error) {
	result := repo.db.
		Preload("Category").
		Preload("SubCategory").
		Preload("Ratings").
		Create(&project).
		First(&project)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create project")
		return nil, result.Error
	}

	return project, nil
}

func (repo *projectRepository) FindByID(projectID uuid.UUID) (*entity.Project, error) {
	var project entity.Project
	result := repo.db.
		Preload("Owner").
		Preload("Category").
		Preload("SubCategory").
		Preload("Ratings").
		Where("id = ?", projectID).
		First(&project)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find project by id")
		return nil, result.Error
	}

	return &project, nil
}

func (repo *projectRepository) FindAllByOwnerID(ownerID uuid.UUID) ([]entity.Project, error) {
	var projects []entity.Project
	result := repo.db.
		Preload("Owner").
		Preload("Category").
		Preload("SubCategory").
		Preload("Ratings").
		Where("owner_id = ?", ownerID).
		Find(&projects)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list projects")
		return nil, result.Error
	}

	return projects, nil
}

func (repo *projectRepository) CreateProjectRating(rating *entity.ProjectRating) (*entity.ProjectRating, error) {
	result := repo.db.
		Create(&rating).
		First(&rating)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create project rating")
		return nil, result.Error
	}

	return rating, nil
}

func (repo *projectRepository) FindProjectRating(userID uuid.UUID, projectID uuid.UUID) (*entity.ProjectRating, error) {
	var rating entity.ProjectRating
	result := repo.db.
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Find(&rating)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find project rating")
		return nil, result.Error
	}

	return &rating, nil
}
