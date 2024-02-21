package repository

import (
	"fund-o/api-server/internal/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *entity.Project) (*entity.Project, error)
	FindAllByOwnerID(ownerID uuid.UUID) ([]entity.Project, error)
}

type projectRepository struct {
	db     *gorm.DB
	logger *log.Entry
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	logger := log.WithFields(log.Fields{
		"module": "project_repository",
	})
	return &projectRepository{db, logger}
}

func (repo *projectRepository) Create(project *entity.Project) (*entity.Project, error) {
	result := repo.db.
		Preload("Category").
		Preload("SubCategory").
		Create(&project).
		First(&project)
	if result.Error != nil {
		repo.logger.Errorf("Failed to create project: %v", result.Error)
		return nil, result.Error
	}

	return project, nil
}

func (repo *projectRepository) FindAllByOwnerID(ownerID uuid.UUID) ([]entity.Project, error) {
	var projects []entity.Project
	result := repo.db.
		Preload("Owner").
		Preload("Category").
		Preload("SubCategory").
		Where("owner_id = ?", ownerID).
		Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}

	return projects, nil
}
