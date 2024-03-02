package handler

import (
	_ "fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandlerOptions struct {
	usecase.UserUseCase
}

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(options *UserHandlerOptions) *UserHandler {
	return &UserHandler{
		userUseCase: options.UserUseCase,
	}
}

// GetMe godoc
// @summary Get current user
// @description Get current user by validating authorization token
// @tags users
// @id GetMe
// @produce json
// @security ApiKeyAuth
// @response 200 {object} handler.ResultResponse[entity.UserDto] "OK"
// @response 401 {object} handler.ErrorResponse "Unauthorized"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	user, err := h.userUseCase.GetUserById(payload.UserID)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, user))
}
