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
