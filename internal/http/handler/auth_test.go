package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fund-o/api-server/cmd/worker"
	"fund-o/api-server/internal/datasource"
	"fund-o/api-server/internal/datasource/driver"
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/mocks"
	"fund-o/api-server/pkg/apperrors"
	"fund-o/api-server/pkg/password"
	"fund-o/api-server/pkg/random"
	"fund-o/api-server/pkg/token"
	"github.com/hibiken/asynq"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerSuite struct {
	suite.Suite
	datasources    datasource.Datasource
	userRepository repository.UserRepository
	handler        *AuthHandler
}

func (s *AuthHandlerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.datasources = datasource.NewDatasourceContext(&datasource.DatasourceConfig{
		SqlDBConfig: driver.SqlDBConfig{
			SQL_HOST:     "localhost",
			SQL_USERNAME: "docker",
			SQL_PASSWORD: "secret",
			SQL_PORT:     5433,
			SQL_DATABASE: "fundo_test",
		},
	})

	s.userRepository = repository.NewUserRepository(s.datasources.GetSqlDB())
	sessionRepository := repository.NewSessionRepository(s.datasources.GetSqlDB())
	imageUploader := mocks.NewMockImageUploader(ctrl)
	userUseCase := usecase.NewUserUseCase(&usecase.UserUseCaseOptions{
		UserRepository: s.userRepository,
		ImageUploader:  imageUploader,
	})
	sessionUseCase := usecase.NewSessionUseCase(&usecase.SessionUseCaseOptions{
		SessionRepository: sessionRepository,
	})

	redisOpts := asynq.RedisClientOpt{
		Addr: "localhost:6380",
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)

	secretKey := "alsypVB6YUpE2HBW4npGoXeArNyqVrqO"

	tokenMaker, err := token.NewJWTMaker(secretKey)
	require.NoError(s.T(), err)

	s.handler = NewAuthHandler(&AuthHandlerOptions{
		UserUseCase:     userUseCase,
		SessionUseCase:  sessionUseCase,
		TaskDistributor: taskDistributor,
		TokenMaker:      tokenMaker,
	})
}

func (s *AuthHandlerSuite) TearDownSuite() {
	err := s.datasources.Close()
	require.NoError(s.T(), err)
}

func (s *AuthHandlerSuite) TestAuthRegisterAPI() {
	validRequestBody := entity.UserCreatePayload{
		Email:                random.NewEmail(),
		Password:             "@Password123",
		PasswordConfirmation: "@Password123",
		Firstname:            "John",
		Lastname:             "Doe",
		BirthDate:            "2002-04-16T00:00:00Z",
		Gender:               "m",
	}

	testCases := []struct {
		name          string
		requestBody   entity.UserCreatePayload
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			requestBody: validRequestBody,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.UserAuthenticateResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusCreated), response.Status)
				require.Equal(t, http.StatusCreated, response.StatusCode)
				require.NotNil(t, response.Result)
				require.NotNil(t, response.Result.SessionID)
				require.NotNil(t, response.Result.AccessToken)

				require.Equal(t, validRequestBody.Email, response.Result.User.Email)
				require.Equal(t, fmt.Sprintf("%s %s", validRequestBody.Firstname, validRequestBody.Lastname), response.Result.User.FullName)
				require.Equal(t, validRequestBody.BirthDate, response.Result.User.BirthDate)
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: entity.UserCreatePayload{},
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
			requestBody: entity.UserCreatePayload{
				Email:                random.NewEmail(),
				Password:             "@Password123",
				PasswordConfirmation: "@Password1234",
				Firstname:            "John",
				Lastname:             "Doe",
				BirthDate:            "2002-04-16T00:00:00Z",
				Gender:               "m",
			},
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
			requestBody: entity.UserCreatePayload{
				Email:                random.NewEmail(),
				Password:             "@Password123",
				PasswordConfirmation: "@Password123",
				Firstname:            "John",
				Lastname:             "Doe",
				BirthDate:            "2002-04-16",
				Gender:               "m",
			},
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
			requestBody: entity.UserCreatePayload{
				Email:                random.NewEmail(),
				Password:             "password",
				PasswordConfirmation: "password",
				Firstname:            "John",
				Lastname:             "Doe",
				BirthDate:            "2002-04-16T00:00:00Z",
				Gender:               "m",
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusInternalServerError), response.Status)
				require.Equal(t, http.StatusInternalServerError, response.StatusCode)
				require.Equal(t, apperrors.ErrHashPassword.Error(), response.Error)
			},
		},
		{
			name:        "EmailAlreadyExists",
			requestBody: validRequestBody,
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

			r.POST("/register", s.handler.Register)

			requestBody, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			request, err := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))

			require.NoError(t, err)

			c.Request = request

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *AuthHandlerSuite) TestAuthLoginAPI() {
	email := random.NewEmail()

	testCases := []struct {
		name          string
		requestBody   entity.UserLoginPayload
		buildStubs    func(t *testing.T, repo repository.UserRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			requestBody: entity.UserLoginPayload{
				Email:    email,
				Password: "@Password123",
			},
			buildStubs: func(t *testing.T, repo repository.UserRepository) {
				hashedPassword, err := password.HashPassword("@Password123")
				require.NoError(t, err)

				birthDate, err := time.Parse(time.RFC3339, "2002-04-16T00:00:00Z")
				require.NoError(t, err)

				_, err = repo.Create(&entity.User{
					Email:          email,
					HashedPassword: hashedPassword,
					Firstname:      "John",
					Lastname:       "Doe",
					DisplayName:    "John Doe",
					BirthDate:      birthDate,
				})
				require.NoError(t, err)
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

				require.Equal(t, email, response.Result.User.Email)
			},
		},
		{
			name:        "InvalidRequestBody",
			requestBody: entity.UserLoginPayload{},
			buildStubs:  func(t *testing.T, repo repository.UserRepository) {},
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

			tc.buildStubs(t, s.userRepository)

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

func TestAuthHandlerSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerSuite))
}
