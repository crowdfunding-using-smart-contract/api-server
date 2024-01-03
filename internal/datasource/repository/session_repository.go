package repository

import (
	"fund-o/api-server/internal/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *entity.Session) (*entity.Session, error)
	FindByID(id uuid.UUID) (*entity.Session, error)
}

type sessionRepository struct {
	db     *gorm.DB
	logger *log.Entry
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	logger := log.WithFields(log.Fields{
		"module": "session_repository",
	})
	return &sessionRepository{db, logger}
}

func (repo *sessionRepository) Create(session *entity.Session) (*entity.Session, error) {
	if result := repo.db.Create(&session); result.Error != nil {
		repo.logger.Errorf("Failed to create session: %v", result.Error)
		return nil, result.Error
	}

	return session, nil
}

func (repo *sessionRepository) FindByID(id uuid.UUID) (*entity.Session, error) {
	var session entity.Session
	if result := repo.db.Where("id = ?", id).First(&session); result.Error != nil {
		repo.logger.Errorf("Failed to find session by id: %v", result.Error)
		return nil, result.Error
	}

	return &session, nil
}
