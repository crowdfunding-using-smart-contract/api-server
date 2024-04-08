package repository

import (
	"fund-o/api-server/internal/entity"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *entity.Message) (*entity.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *entity.Message) (*entity.Message, error) {
	result := r.db.
		Preload("Author").
		Create(&message).
		First(&message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}
