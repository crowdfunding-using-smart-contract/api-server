package repository

import (
	"fund-o/api-server/internal/entity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type VerifyEmailRepository interface {
	Create(verifyEmail *entity.VerifyEmail) (*entity.VerifyEmail, error)
	FindByID(id string) (*entity.VerifyEmail, error)
	UpdateByID(id string, verifyEmail *entity.VerifyEmail) (*entity.VerifyEmail, error)
}

type verifyEmailRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewVerifyEmailRepository(db *gorm.DB) VerifyEmailRepository {
	logger := log.With().Str("module", "verify_email_repository").Logger()
	return &verifyEmailRepository{db, logger}
}

func (repo *verifyEmailRepository) Create(verifyEmail *entity.VerifyEmail) (*entity.VerifyEmail, error) {
	if result := repo.db.Create(&verifyEmail); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create verify email: " + verifyEmail.Email)
		return nil, result.Error
	}
	return verifyEmail, nil
}

func (repo *verifyEmailRepository) FindByID(id string) (*entity.VerifyEmail, error) {
	var ve entity.VerifyEmail
	if result := repo.db.Where("id = ?", id).First(&ve); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find verify email by id: " + id)
		return nil, result.Error
	}
	return &ve, nil
}

func (repo *verifyEmailRepository) UpdateByID(id string, verifyEmail *entity.VerifyEmail) (*entity.VerifyEmail, error) {
	if result := repo.db.Model(&entity.VerifyEmail{}).Where("id = ?", id).Updates(verifyEmail); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to update verify email by id: " + id)
		return nil, result.Error
	}
	return verifyEmail, nil
}
