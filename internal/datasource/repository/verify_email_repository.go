package repository

import (
	"fund-o/api-server/internal/entity"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VerifyEmailRepository interface {
	Create(verifyEmail *entity.VerifyEmail) (*entity.VerifyEmail, error)
}

type verifyEmailRepository struct {
	db     *gorm.DB
	logger *log.Entry
}

func NewVerifyEmailRepository(db *gorm.DB) VerifyEmailRepository {
	logger := log.WithFields(log.Fields{
		"module": "verify_email_repository",
	})
	return &verifyEmailRepository{db, logger}
}

func (repo *verifyEmailRepository) Create(verifyEmail *entity.VerifyEmail) (*entity.VerifyEmail, error) {
	if result := repo.db.Create(&verifyEmail); result.Error != nil {
		repo.logger.Errorf("Failed to create verify email: %v", result.Error)
		return nil, result.Error
	}
	return verifyEmail, nil
}
