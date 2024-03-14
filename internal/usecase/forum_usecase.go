package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/pagination"
	"github.com/google/uuid"
)

type ForumUseCase interface {
	ListForums(paginateOptions pagination.PaginateOptions) pagination.PaginateResult[entity.PostDto]
	CreatePost(payload *entity.PostCreatePayload) (*entity.PostDto, error)
	GetPostByID(id string) (*entity.PostDto, error)
	CreateCommentByForumID(forumID string, comment *entity.CommentCreatePayload) (*entity.CommentDto, error)
	CreateReplyByCommentID(commentID string, payload *entity.ReplyCreatePayload) (*entity.ReplyDto, error)
}

type forumUseCase struct {
	forumRepository repository.ForumRepository
}

type ForumUseCaseOptions struct {
	repository.ForumRepository
}

func NewForumUseCase(options *ForumUseCaseOptions) ForumUseCase {
	return &forumUseCase{
		forumRepository: options.ForumRepository,
	}
}

func (uc *forumUseCase) ListForums(paginateOptions pagination.PaginateOptions) pagination.PaginateResult[entity.PostDto] {
	result := pagination.MakePaginateResult(pagination.MakePaginateContextParameters[entity.PostDto]{
		PaginateOptions: paginateOptions,
		CountDocuments: func() int64 {
			return uc.forumRepository.CountPost()
		},
		FindDocuments: func(findOptions pagination.PaginateFindOptions) []entity.PostDto {
			documents := uc.forumRepository.ListPosts(findOptions)

			forumDtos := make([]entity.PostDto, 0, len(documents))
			for _, document := range documents {
				forumDtos = append(forumDtos, *document.ToPostDto())
			}

			return forumDtos
		},
	})

	return result
}

func (uc *forumUseCase) CreatePost(payload *entity.PostCreatePayload) (*entity.PostDto, error) {
	forum, err := uc.forumRepository.CreatePost(&entity.Post{
		Title:    payload.Title,
		Content:  payload.Content,
		AuthorID: uuid.MustParse(payload.AuthorID),
	})
	if err != nil {
		return nil, err
	}

	return forum.ToPostDto(), nil
}

func (uc *forumUseCase) GetPostByID(id string) (*entity.PostDto, error) {
	forum, err := uc.forumRepository.FindPostByID(uuid.MustParse(id))
	if err != nil {
		return nil, err
	}

	return forum.ToPostDto(), nil
}

func (uc *forumUseCase) CreateCommentByForumID(forumID string, payload *entity.CommentCreatePayload) (*entity.CommentDto, error) {
	comment, err := uc.forumRepository.CreateComment(&entity.Comment{
		Content:  payload.Content,
		AuthorID: uuid.MustParse(payload.AuthorID),
		ForumID:  uuid.MustParse(forumID),
	})
	if err != nil {
		return nil, err
	}

	return comment.ToCommentDto(), nil
}

func (uc *forumUseCase) CreateReplyByCommentID(commentID string, payload *entity.ReplyCreatePayload) (*entity.ReplyDto, error) {
	reply, err := uc.forumRepository.CreateReply(&entity.Reply{
		Content:   payload.Content,
		AuthorID:  uuid.MustParse(payload.AuthorID),
		CommentID: uuid.MustParse(commentID),
	})
	if err != nil {
		return nil, err
	}

	return reply.ToReplyDto(), nil
}
