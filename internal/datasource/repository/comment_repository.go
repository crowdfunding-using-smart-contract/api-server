package repository

import (
	"fund-o/api-server/internal/entity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *entity.Comment) (*entity.Comment, error)
	CreateReply(reply *entity.Reply) (*entity.Reply, error)
}

type commentRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	logger := log.With().Str("module", "comment_repository").Logger()
	return &commentRepository{db, logger}
}

func (repo *commentRepository) Create(comment *entity.Comment) (*entity.Comment, error) {
	result := repo.db.
		Preload("Author").
		Create(&comment).
		First(&comment)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create comment")
		return nil, result.Error
	}

	return comment, nil
}

func (repo *commentRepository) CreateReply(reply *entity.Reply) (*entity.Reply, error) {
	result := repo.db.
		Create(&reply).
		First(&reply)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create reply")
		return nil, result.Error
	}

	return reply, nil
}
