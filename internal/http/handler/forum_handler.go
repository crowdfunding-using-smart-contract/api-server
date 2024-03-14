package handler

import (
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/pagination"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ForumHandler struct {
	forumUseCase usecase.ForumUseCase
}

type ForumHandlerOptions struct {
	usecase.ForumUseCase
}

func NewForumHandler(options *ForumHandlerOptions) *ForumHandler {
	return &ForumHandler{
		forumUseCase: options.ForumUseCase,
	}
}

// ListPosts godoc
// @summary List Posts
// @description List posts
// @tags forums
// @id ListPosts
// @accept json
// @produce json
// @param page query int false "number of page"
// @param size query int false "size of data per page"
// @success 200 {object} handler.ResultResponse[pagination.PaginateResult[entity.PostDto]] "OK"
// @failure 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /posts [get]
func (h *ForumHandler) ListPosts(c *gin.Context) {
	var paginateOptions pagination.PaginateOptions
	if err := c.ShouldBindQuery(&paginateOptions); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error list forums: %v", err.Error())))
	}

	forums := h.forumUseCase.ListForums(paginateOptions)
	c.JSON(makeHttpResponse(http.StatusOK, forums))
}

// CreatePost godoc
// @summary Create Post
// @description Create post
// @tags forums
// @id CreatePost
// @accept json
// @produce json
// @security ApiKeyAuth
// @param payload body entity.PostCreatePayload true "post payload"
// @success 201 {object} handler.ResultResponse[entity.PostDto]
// @failure 400 {object} handler.ErrorResponse
// @failure 500 {object} handler.ErrorResponse
// @router /posts [post]
func (h *ForumHandler) CreatePost(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID
	var req entity.PostCreatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error create forum: %v", err.Error())))
	}

	req.AuthorID = userID

	forumDto, err := h.forumUseCase.CreatePost(&req)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error create forum: %v", err.Error())))
	}

	c.JSON(makeHttpResponse(http.StatusCreated, forumDto))
}

// GetPostByID godoc
// @summary Get Post by ID
// @description Get post by id
// @tags forums
// @id GetPostByID
// @accept json
// @produce json
// @param id path string true "post id to get"
// @success 200 {object} handler.ResultResponse[entity.PostDto]
// @failure 500 {object} handler.ErrorResponse
// @router /posts/{id} [get]
func (h *ForumHandler) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	forumDto, err := h.forumUseCase.GetPostByID(id)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error get forum by id: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, forumDto))
}

// CreateComment godoc
// @summary Create Comment
// @description Create comment for forum
// @tags forums
// @id CreateComment
// @accept json
// @produce json
// @security ApiKeyAuth
// @param id path string true "forum id to comment"
// @param payload body entity.CommentCreatePayload true "comment payload"
// @success 201 {object} handler.ResultResponse[entity.CommentDto]
// @failure 400 {object} handler.ErrorResponse
// @failure 500 {object} handler.ErrorResponse
// @router /posts/{id}/comments [post]
func (h *ForumHandler) CreateComment(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID
	forumID := c.Param("id")
	var req entity.CommentCreatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error create comment: %v", err.Error())))
	}

	req.AuthorID = userID

	commentDto, err := h.forumUseCase.CreateCommentByForumID(forumID, &req)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error create comment: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusCreated, commentDto))
}