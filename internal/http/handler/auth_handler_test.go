package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/mocks"
	"fund-o/api-server/pkg/apperrors"
	"fund-o/api-server/pkg/random"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthHandlerSuite struct {
	suite.Suite
	userRepository    *mocks.MockUserRepository
	sessionRepository *mocks.MockSessionRepository
	taskDistributor   *mocks.MockTaskDistributor
	handler           *AuthHandler
}

func (s *AuthHandlerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.userRepository = mocks.NewMockUserRepository(ctrl)
	s.sessionRepository = mocks.NewMockSessionRepository(ctrl)
	imageUploader := mocks.NewMockImageUploader(ctrl)
	userUseCase := usecase.NewUserUseCase(&usecase.UserUseCaseOptions{
		UserRepository: s.userRepository,
		ImageUploader:  imageUploader,
	})
	sessionUseCase := usecase.NewSessionUseCase(&usecase.SessionUseCaseOptions{
		SessionRepository: s.sessionRepository,
	})

	s.taskDistributor = mocks.NewMockTaskDistributor(ctrl)

	secretKey := "alsypVB6YUpE2HBW4npGoXeArNyqVrqO"

	tokenMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(s.T(), err)

	s.handler = NewAuthHandler(&AuthHandlerOptions{
		UserUseCase:     userUseCase,
		SessionUseCase:  sessionUseCase,
		TaskDistributor: s.taskDistributor,
		TokenMaker:      tokenMaker,
	})
}

func (s *AuthHandlerSuite) TestAuthRegisterAPI() {
	user := randomUser(s.T())

	testCases := []struct {
		name          string
		requestBody   gin.H
		buildStubs    func()
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"email":                 user.Email,
				"password":              "@Password123",
				"password_confirmation": "@Password123",
				"firstname":             user.Firstname,
				"lastname":              user.Lastname,
				"birthdate":             "2000-01-01T00:00:00Z",
				"gender":                "m",
			},
			buildStubs: func() {
				s.userRepository.EXPECT().
					Create(gomock.Any()).
					Times(1).
					Return(&user, nil)

				s.sessionRepository.EXPECT().
					Create(gomock.Any()).
					Times(1).
					Return(&entity.Session{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.UserAuthenticateResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusCreated), response.Status)
				require.Equal(t, http.StatusCreated, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, user.Email, response.Result.User.Email)
				require.Equal(t, fmt.Sprintf("%s %s", user.Firstname, user.Lastname), response.Result.User.FullName)
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: gin.H{},
			buildStubs:  func() {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "PasswordDoesNotMatch",
			requestBody: gin.H{
				"email":                 random.NewEmail(),
				"password":              "@Password123",
				"password_confirmation": "@Password1234",
				"firstname":             "John",
				"lastname":              "Doe",
				"birthdate":             "2002-04-16T00:00:00Z",
				"gender":                "m",
			},
			buildStubs: func() {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
				require.Equal(t, apperrors.ErrPasswordAndConfirmationNotMatch.Error(), response.Error)
			},
		},
		{
			name: "InvalidBirthDate",
			requestBody: gin.H{
				"email":                 random.NewEmail(),
				"password":              "@Password123",
				"password_confirmation": "@Password123",
				"firstname":             "John",
				"lastname":              "Doe",
				"birthdate":             "2002-04-16",
				"gender":                "m",
			},
			buildStubs: func() {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
				require.Equal(t, apperrors.ErrInvalidBirthDateFormat.Error(), response.Error)
			},
		},
		{
			name: "HashPasswordError",
			requestBody: gin.H{
				"email":                 random.NewEmail(),
				"password":              "password",
				"password_confirmation": "password",
				"firstname":             "John",
				"lastname":              "Doe",
				"birthdate":             "2000-01-01T00:00:00Z",
				"gender":                "m",
			},
			buildStubs: func() {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
				require.Equal(t, apperrors.ErrHashPassword.Error(), response.Error)
			},
		},
		{
			name: "EmailAlreadyExists",
			requestBody: gin.H{
				"email":                 user.Email,
				"password":              "@Password123",
				"password_confirmation": "@Password123",
				"firstname":             user.Firstname,
				"lastname":              user.Lastname,
				"birthdate":             "2000-01-01T00:00:00Z",
				"gender":                "m",
			},
			buildStubs: func() {
				s.userRepository.EXPECT().
					Create(gomock.Any()).
					Times(1).
					Return(nil, gorm.ErrDuplicatedKey)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusInternalServerError), response.Status)
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			tc.buildStubs()

			r.POST("/register", s.handler.Register)

			requestBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(requestBody))
			require.NoError(t, err)

			c.Request = request

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *AuthHandlerSuite) TestAuthLoginAPI() {
	user := randomUser(s.T())

	testCases := []struct {
		name          string
		requestBody   entity.UserLoginPayload
		buildStubs    func(userRepo *mocks.MockUserRepository, sessionRepo *mocks.MockSessionRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: entity.UserLoginPayload{
				Email:    user.Email,
				Password: "@Password123",
			},
			buildStubs: func(userRepo *mocks.MockUserRepository, sessionRepo *mocks.MockSessionRepository) {
				userRepo.EXPECT().
					FindByEmail(user.Email).
					Times(1).
					Return(&user, nil)

				sessionRepo.EXPECT().
					Create(gomock.Any()).
					Times(1).
					Return(&entity.Session{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.UserAuthenticateResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)

				require.NotNil(t, response.Result)
				require.NotNil(t, response.Result.SessionID)
				require.NotNil(t, response.Result.AccessToken)
				require.NotNil(t, response.Result.User, "User should not be nil")

				require.Equal(t, user.Email, response.Result.User.Email)
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: entity.UserLoginPayload{},
			buildStubs:  func(userRepo *mocks.MockUserRepository, sessionRepo *mocks.MockSessionRepository) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "UserNotFound",
			requestBody: entity.UserLoginPayload{
				Email:    user.Email,
				Password: "@Password123",
			},
			buildStubs: func(userRepo *mocks.MockUserRepository, sessionRepo *mocks.MockSessionRepository) {
				userRepo.EXPECT().
					FindByEmail(user.Email).
					Times(1).
					Return(nil, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusNotFound), response.Status)
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			tc.buildStubs(s.userRepository, s.sessionRepository)

			r.POST("/login", s.handler.Login)

			requestBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
			require.NoError(t, err)

			c.Request = request

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *AuthHandlerSuite) TestSendVerifyEmailAPI() {
	email := random.NewEmail()

	testCases := []struct {
		name          string
		requestBody   gin.H
		buildStubs    func()
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: gin.H{
				"email": email,
			},
			buildStubs: func() {
				s.userRepository.EXPECT().
					FindByEmail(email).
					Times(1).
					Return(&entity.User{}, nil)

				s.taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response MessageResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: gin.H{},
			buildStubs:  func() {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "UserNotFound",
			requestBody: gin.H{
				"email": email,
			},
			buildStubs: func() {
				s.userRepository.EXPECT().
					FindByEmail(email).
					Times(1).
					Return(nil, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			tc.buildStubs()

			r.POST("/send-verify-email", s.handler.SendVerifyEmail)

			requestBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/send-verify-email", bytes.NewReader(requestBody))
			require.NoError(t, err)

			c.Request = request

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}
