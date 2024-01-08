package repository

import (
	"fund-o/api-server/internal/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *entity.Project) (*entity.Project, error)
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
	if result := repo.db.Create(&project); result.Error != nil {
		repo.logger.Errorf("Failed to create project: %v", result.Error)
		return nil, result.Error
	}

	return project, nil
}
