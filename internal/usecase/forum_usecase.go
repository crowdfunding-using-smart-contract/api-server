package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/pagination"
	"github.com/google/uuid"
)

type ForumUseCase interface {
	ListForums(paginateOptions pagination.PaginateOptions) pagination.PaginateResult[entity.ForumDto]
	CreateForum(payload *entity.ForumCreatePayload) (*entity.ForumDto, error)
	GetForumByID(id string) (*entity.ForumDto, error)
	CreateCommentByForumID(forumID string, comment *entity.CommentCreatePayload) (*entity.CommentDto, error)
}

type forumUseCase struct {
	forumRepository   repository.ForumRepository
	commentRepository repository.CommentRepository
}

type ForumUseCaseOptions struct {
	repository.ForumRepository
	repository.CommentRepository
}

func NewForumUseCase(options *ForumUseCaseOptions) ForumUseCase {
	return &forumUseCase{
		forumRepository:   options.ForumRepository,
		commentRepository: options.CommentRepository,
	}
}

func (uc *forumUseCase) ListForums(paginateOptions pagination.PaginateOptions) pagination.PaginateResult[entity.ForumDto] {
	result := pagination.MakePaginateResult(pagination.MakePaginateContextParameters[entity.ForumDto]{
		PaginateOptions: paginateOptions,
		CountDocuments: func() int64 {
			return uc.forumRepository.Count()
		},
		FindDocuments: func(findOptions pagination.PaginateFindOptions) []entity.ForumDto {
			documents := uc.forumRepository.List(findOptions)

			forumDtos := make([]entity.ForumDto, 0, len(documents))
			for _, document := range documents {
				forumDtos = append(forumDtos, *document.ToForumDto())
			}

			return forumDtos
		},
	})

	return result
}

func (uc *forumUseCase) CreateForum(payload *entity.ForumCreatePayload) (*entity.ForumDto, error) {
	forum, err := uc.forumRepository.Create(&entity.Forum{
		Title:    payload.Title,
		Content:  payload.Content,
		AuthorID: uuid.MustParse(payload.AuthorID),
	})
	if err != nil {
		return nil, err
	}

	return forum.ToForumDto(), nil
}

func (uc *forumUseCase) GetForumByID(id string) (*entity.ForumDto, error) {
	forum, err := uc.forumRepository.FindByID(uuid.MustParse(id))
	if err != nil {
		return nil, err
	}

	return forum.ToForumDto(), nil
}

func (uc *forumUseCase) CreateCommentByForumID(forumID string, payload *entity.CommentCreatePayload) (*entity.CommentDto, error) {
	comment, err := uc.commentRepository.Create(&entity.Comment{
		Content:  payload.Content,
		AuthorID: uuid.MustParse(payload.AuthorID),
		ForumID:  uuid.MustParse(forumID),
	})
	if err != nil {
		return nil, err
	}

	return comment.ToCommentDto(), nil
}
