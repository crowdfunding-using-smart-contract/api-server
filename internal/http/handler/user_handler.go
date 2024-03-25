package handler

import (
	"errors"
	"fmt"
	"fund-o/api-server/internal/entity"
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

// UpdateUser godoc
// @summary Update user
// @description Update user by id
// @tags users
// @id UpdateUser
// @accept multipart/form-data
// @produce json
// @security ApiKeyAuth
// @param id path string true "User ID"
// @param image formData file true "User profile image"
// @response 200 {object} handler.ResultResponse[entity.UserDto] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 401 {object} handler.ErrorResponse "Unauthorized"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID

	if id != userID {
		err := errors.New("unauthorized: user can only update their own profile")
		c.JSON(makeHttpErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	var payload entity.UserUpdatePayload
	if err := c.ShouldBind(&payload); err != nil {
		fmt.Println("error", err.Error())
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	user, err := h.userUseCase.UpdateUserByID(id, &payload)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, user))
}
