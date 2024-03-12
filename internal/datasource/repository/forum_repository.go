package repository

import (
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type forumRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

type ForumRepository interface {
	List(findOptions pagination.PaginateFindOptions) []entity.Forum
	Count() int64
	Create(forum *entity.Forum) (*entity.Forum, error)
	FindByID(id uuid.UUID) (*entity.Forum, error)
	FindAllByAuthorID(authorID uuid.UUID) ([]entity.Forum, error)
}

func NewForumRepository(db *gorm.DB) ForumRepository {
	logger := log.With().Str("module", "forum_repository").Logger()
	return &forumRepository{db, logger}
}

func (repo *forumRepository) List(findOptions pagination.PaginateFindOptions) (forums []entity.Forum) {
	if result := repo.db.Limit(findOptions.Limit).Offset(findOptions.Skip).Find(&forums); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list forums")
		return
	}

	return forums
}

func (repo *forumRepository) Count() int64 {
	var count int64
	if result := repo.db.Model(&entity.Forum{}).Count(&count); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to count forums")
		return 0
	}

	return count
}

func (repo *forumRepository) Create(forum *entity.Forum) (*entity.Forum, error) {
	result := repo.db.
		Preload("Author").
		Create(&forum).
		First(&forum)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create forum")
		return nil, result.Error
	}

	return forum, nil
}

func (repo *forumRepository) FindByID(id uuid.UUID) (*entity.Forum, error) {
	var forum entity.Forum
	result := repo.db.
		Preload("Author").
		Preload("Comments").
		Where("id = ?", id).
		First(&forum)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find forum")
		return nil, result.Error
	}

	return &forum, nil
}

func (repo *forumRepository) FindAllByAuthorID(authorID uuid.UUID) ([]entity.Forum, error) {
	var forums []entity.Forum
	result := repo.db.
		Preload("Author").
		Where("author_id = ?", authorID).
		Find(&forums)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list forums")
		return nil, result.Error
	}

	return forums, nil
}
