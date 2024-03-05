package repository

import (
	"fund-o/api-server/internal/entity"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindById(id uuid.UUID) (*entity.User, error)
	UpdateByID(id uuid.UUID, user *entity.User) (*entity.User, error)
}

type userRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewUserRepository(db *gorm.DB) UserRepository {
	logger := log.With().Str("module", "user_repository").Logger()
	return &userRepository{db, logger}
}

func (repo *userRepository) Create(user *entity.User) (*entity.User, error) {
	if result := repo.db.Create(&user); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create user")
		return nil, result.Error
	}

	return user, nil
}

func (repo *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if result := repo.db.Where("email = ?", email).First(&user); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find user by email: " + email)
		return nil, result.Error
	}

	return &user, nil
}

func (repo *userRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if result := repo.db.First(&user, id); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find user by id: " + id.String())
		return nil, result.Error
	}

	return &user, nil
}

func (repo *userRepository) UpdateByID(id uuid.UUID, user *entity.User) (*entity.User, error) {
	if result := repo.db.Model(&entity.User{}).Where("id = ?", id).Updates(&user).First(&user); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to update user by id: " + id.String())
		return nil, result.Error
	}

	return user, nil
}
