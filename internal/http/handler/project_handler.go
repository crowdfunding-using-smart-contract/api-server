package handler

import (
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectUseCase         usecase.ProjectUseCase
	projectCategoryUseCase usecase.ProjectCategoryUseCase
	userUseCase            usecase.UserUseCase
}

type ProjectHandlerOptions struct {
	usecase.ProjectUseCase
	usecase.ProjectCategoryUseCase
	usecase.UserUseCase
}

func NewProjectHandler(options *ProjectHandlerOptions) *ProjectHandler {
	return &ProjectHandler{
		projectUseCase:         options.ProjectUseCase,
		projectCategoryUseCase: options.ProjectCategoryUseCase,
		userUseCase:            options.UserUseCase,
	}
}

// CreateProject godoc
// @summary Create Project
// @description Create project with required data
// @tags projects
// @id CreateProject
// @accpet json
// @produce json
// @security ApiKeyAuth
// @param Project body entity.ProjectCreatePayload true "Project data to be created"
// @response 200 {object} handler.ResultResponse[entity.ProjectDto] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID

	var req entity.ProjectCreatePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error create project: %v", err.Error())))
		return
	}

	user, err := h.userUseCase.GetUserById(userID)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error create project: %v", err.Error())))
		return
	}

	req.OwnerID = user.ID

	projectDto, err := h.projectUseCase.CreateProject(&req)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error create project: %v", err.Error())))
		return
	}

	projectDto.Owner = user

	c.JSON(makeHttpResponse(http.StatusCreated, projectDto))
}

// GetOwnProjects godoc
// @summary Get own Projects
// @description Get own projects with authenticate creator
// @tags projects
// @id GetOwnProjects
// @accpet json
// @produce json
// @security ApiKeyAuth
// @response 200 {object} handler.ResultResponse[[]entity.ProjectDto] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /projects/own [get]
func (h *ProjectHandler) GetOwnProjects(c *gin.Context) {
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID

	projectDtos, err := h.projectUseCase.GetProjectsByOwnerID(userID)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error create project: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, projectDtos))
}

// ListProjectCategories godoc
// @summary List Project Categories
// @description List project categories for selection
// @tags projects
// @id ListProjectCategories
// @produce json
// @response 200 {object} handler.ResultResponse[[]entity.ProjectCategoryDto]
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /projects/categories [get]
func (h *ProjectHandler) ListProjectCategories(c *gin.Context) {
	categories, err := h.projectCategoryUseCase.ListProjectCategories()
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error list project categories: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, categories))
}
