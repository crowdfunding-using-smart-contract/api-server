package handler

import (
	"errors"
	"fmt"
	"fund-o/api-server/cmd/worker"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/pkg/token"
	"github.com/hibiken/asynq"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandlerOptions struct {
	usecase.UserUseCase
	usecase.SessionUseCase
	TokenMaker token.Maker
	worker.TaskDistributor
}

type AuthHandler struct {
	userUseCase     usecase.UserUseCase
	sessionUseCase  usecase.SessionUseCase
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewAuthHandler(options *AuthHandlerOptions) *AuthHandler {
	return &AuthHandler{
		userUseCase:     options.UserUseCase,
		sessionUseCase:  options.SessionUseCase,
		tokenMaker:      options.TokenMaker,
		taskDistributor: options.TaskDistributor,
	}
}

// Register godoc
// @summary Register User
// @description Create user with specific user data and role
// @tags auth
// @id Register
// @accpet json
// @produce json
// @param User body entity.UserCreatePayload true "User data to be created"
// @response 200 {object} handler.ResultResponse[entity.UserDto] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var user entity.UserCreatePayload
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error register user: %v", err.Error())))
		return
	}

	userDto, err := h.userUseCase.CreateUser(&user)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error register user: %v", err.Error())))
		return
	}

	taskPayload := &worker.PayloadSendVerifyEmail{
		Email: userDto.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	h.taskDistributor.DistributeTaskSendVerifyEmail(c, taskPayload, opts...)

	c.JSON(makeHttpResponse(http.StatusCreated, userDto))
}

// Login godoc
// @summary Authenticate User
// @description Authenticate user with email and password
// @tags auth
// @id Login
// @accpet json
// @produce json
// @param User body entity.UserLoginPayload true "User data to be authenticated"
// @response 200 {object} handler.ResultResponse[entity.UserLoginResponse] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.UserLoginPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error authenticate user: %v", err.Error())))
		return
	}

	user, err := h.userUseCase.AuthenticateUser(&req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(makeHttpErrorResponse(http.StatusNotFound, fmt.Sprintf("error authenticate user: %v", err.Error())))
			return
		}
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error authenticate user: %v", err.Error())))
		return
	}

	accessToken, accessTokenPayload, err := h.tokenMaker.CreateToken(user.ID, 15*time.Minute)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error authenticate user: %v", err.Error())))
		return
	}

	refreshToken, refreshTokenPayload, err := h.tokenMaker.CreateToken(user.ID, 24*time.Hour)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error authenticate user: %v", err.Error())))
		return
	}

	session, err := h.sessionUseCase.CreateSession(&entity.SessionCreatePayload{
		ID:           refreshTokenPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIP:     c.ClientIP(),
		ExpiredAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error authenticate user: %v", err.Error())))
		return
	}

	response := entity.UserLoginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiredAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
		User:                  user,
	}

	c.JSON(makeHttpResponse(http.StatusOK, response))
}

type RenewAccessTokenPayload struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

// RenewAccessToken godoc
// @summary Renew Access Token
// @description Renew access token with refresh token
// @tags auth
// @id RenewAccessToken
// @accpet json
// @produce json
// @param User body handler.RenewAccessTokenPayload true "Refresh token to be renewed"
// @response 200 {object} handler.ResultResponse[handler.RenewAccessTokenResponse] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /auth/renew-token [post]
func (h *AuthHandler) RenewAccessToken(c *gin.Context) {
	var req RenewAccessTokenPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	refreshTokenPayload, err := h.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
	}

	session, err := h.sessionUseCase.GetSessionByID(refreshTokenPayload.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(makeHttpErrorResponse(http.StatusNotFound, err.Error()))
			return
		}

		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if session.UserID != refreshTokenPayload.UserID {
		err := fmt.Errorf("mismatch session token")
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("session expired")
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	accessToken, accessTokenPayload, err := h.tokenMaker.CreateToken(refreshTokenPayload.UserID, 15*time.Minute)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	response := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiredAt: accessTokenPayload.ExpiredAt,
	}

	c.JSON(makeHttpResponse(http.StatusOK, response))
}
