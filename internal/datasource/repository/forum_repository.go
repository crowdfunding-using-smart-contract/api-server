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
	ListPosts(findOptions pagination.PaginateFindOptions) []entity.Post
	CountPost() int64
	CreatePost(forum *entity.Post) (*entity.Post, error)
	FindPostByID(id uuid.UUID) (*entity.Post, error)
	FindAllPostsByAuthorID(authorID uuid.UUID) ([]entity.Post, error)
	CreateComment(comment *entity.Comment) (*entity.Comment, error)
	CreateReply(reply *entity.Reply) (*entity.Reply, error)
}

func NewForumRepository(db *gorm.DB) ForumRepository {
	logger := log.With().Str("module", "forum_repository").Logger()
	return &forumRepository{db, logger}
}

func (repo *forumRepository) ListPosts(findOptions pagination.PaginateFindOptions) (forums []entity.Post) {
	if result := repo.db.Limit(findOptions.Limit).Offset(findOptions.Skip).Find(&forums); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list posts")
		return
	}

	return forums
}

func (repo *forumRepository) CountPost() int64 {
	var count int64
	if result := repo.db.Model(&entity.Post{}).Count(&count); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to count posts")
		return 0
	}

	return count
}

func (repo *forumRepository) CreatePost(forum *entity.Post) (*entity.Post, error) {
	result := repo.db.
		Preload("Author").
		Create(&forum).
		First(&forum)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create post")
		return nil, result.Error
	}

	return forum, nil
}

func (repo *forumRepository) FindPostByID(id uuid.UUID) (*entity.Post, error) {
	var forum entity.Post
	result := repo.db.
		Preload("Author").
		Preload("Comments").
		Preload("Comments.Author").
		Where("id = ?", id).
		First(&forum)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find post by id: " + id.String())
		return nil, result.Error
	}

	return &forum, nil
}

func (repo *forumRepository) FindAllPostsByAuthorID(authorID uuid.UUID) ([]entity.Post, error) {
	var forums []entity.Post
	result := repo.db.
		Preload("Author").
		Where("author_id = ?", authorID).
		Find(&forums)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list posts by author id: " + authorID.String())
		return nil, result.Error
	}

	return forums, nil
}

func (repo *forumRepository) CreateComment(comment *entity.Comment) (*entity.Comment, error) {
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

func (repo *forumRepository) CreateReply(reply *entity.Reply) (*entity.Reply, error) {
	result := repo.db.
		Preload("Author").
		Create(&reply).
		First(&reply)
	if result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create reply")
		return nil, result.Error
	}

	return reply, nil
}
