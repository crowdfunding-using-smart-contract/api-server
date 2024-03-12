package repository

import (
	"fund-o/api-server/internal/entity"
	"github.com/rs/zerolog"

	"github.com/google/uuid"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *entity.Session) (*entity.Session, error)
	FindByID(id uuid.UUID) (*entity.Session, error)
}

type sessionRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	logger := log.With().Str("module", "session_repository").Logger()
	return &sessionRepository{db, logger}
}

func (repo *sessionRepository) Create(session *entity.Session) (*entity.Session, error) {
	if result := repo.db.Create(&session); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create session")
		return nil, result.Error
	}

	return session, nil
}

func (repo *sessionRepository) FindByID(id uuid.UUID) (*entity.Session, error) {
	var session entity.Session
	if result := repo.db.Where("id = ?", id).First(&session); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find session by id: " + id.String())
		return nil, result.Error
	}

	return &session, nil
}
