package repository

import (
	"fund-o/api-server/internal/entity"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	logger *log.Entry
}

func NewUserRepository(db *gorm.DB) UserRepository {
	logger := log.WithFields(log.Fields{
		"module": "user_repository",
	})
	return &userRepository{db, logger}
}

func (repo *userRepository) Create(user *entity.User) (*entity.User, error) {
	if result := repo.db.Create(&user); result.Error != nil {
		repo.logger.Errorf("Failed to create user: %v", result.Error)
		return nil, result.Error
	}

	return user, nil
}

func (repo *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if result := repo.db.Where("email = ?", email).First(&user); result.Error != nil {
		repo.logger.Errorf("Failed to find user by email: %v", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (repo *userRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if result := repo.db.First(&user, id); result.Error != nil {
		repo.logger.Errorf("Failed to find user by id: %v", result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (repo *userRepository) UpdateByID(id uuid.UUID, user *entity.User) (*entity.User, error) {
	if result := repo.db.Model(&entity.User{}).Where("id = ?", id).Updates(&user).First(&user); result.Error != nil {
		repo.logger.Errorf("Failed to update user by id: %v", result.Error)
		return nil, result.Error
	}

	return user, nil
}
