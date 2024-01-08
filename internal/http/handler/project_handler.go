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
	projectUsecase usecase.ProjectUsecase
}

type ProjectHandlerOptions struct {
	usecase.ProjectUsecase
}

func NewProjectHandler(options *ProjectHandlerOptions) *ProjectHandler {
	return &ProjectHandler{
		projectUsecase: options.ProjectUsecase,
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

	req.OwnerID = userID

	projectDto, err := h.projectUsecase.CreateProject(&req)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error create project: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusCreated, projectDto))
}
