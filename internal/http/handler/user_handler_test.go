package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/mocks"
	"fund-o/api-server/pkg/token"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserTestSuite struct {
	suite.Suite
	tokenMaker     token.Maker
	userRepository *mocks.MockUserRepository
	handler        *UserHandler
}

func (s *UserTestSuite) SetupSuite() {
	var err error
	secretKey := "alsypVB6YUpE2HBW4npGoXeArNyqVrqO"

	s.tokenMaker, err = token.NewJWTMaker(secretKey)
	require.NoError(s.T(), err)

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.userRepository = mocks.NewMockUserRepository(ctrl)
	useUseCase := usecase.NewUserUseCase(&usecase.UserUseCaseOptions{
		UserRepository: s.userRepository,
	})
	s.handler = NewUserHandler(&UserHandlerOptions{
		UserUseCase: useUseCase,
	})
}

func (s *UserTestSuite) TestGetUserAPI() {
	user := randomUser(s.T())

	testCases := []struct {
		name          string
		buildStubs    func(store *mocks.MockUserRepository)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					FindById(gomock.Eq(user.ID)).
					Times(1).
					Return(&user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.UserDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, user.ID.String(), response.Result.ID)
				require.Equal(t, user.Email, response.Result.Email)
				require.Equal(t, fmt.Sprintf("%s %s", user.Firstname, user.Lastname), response.Result.FullName)
				require.Equal(t, user.DisplayName, response.Result.DisplayName)
			},
		},
		{
			name:       "Invalid User ID",
			buildStubs: func(repo *mocks.MockUserRepository) {},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, "invalid-id", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name: "Not Found",
			buildStubs: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					FindById(gomock.Eq(user.ID)).
					Times(1).
					Return(nil, gorm.ErrRecordNotFound)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusNotFound), response.Status)
				require.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(repo *mocks.MockUserRepository) {
				repo.EXPECT().
					FindById(gomock.Eq(user.ID)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
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
			tc.buildStubs(s.userRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/users/me", middleware.AuthMiddleware(s.tokenMaker), s.handler.GetMe)

			request, err := http.NewRequest(http.MethodGet, "/users/me", nil)
			require.NoError(t, err)

			c.Request = request

			tc.setupAuth(t, c.Request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *UserTestSuite) TestUpdateUserAPI() {
	user := randomUser(s.T())
	newDisplayName := "iamfinethankyouandryu"

	testCases := []struct {
		name          string
		userID        string
		body          gin.H
		buildStubs    func(repo *mocks.MockUserRepository)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID.String(),
			body: gin.H{
				"display_name": newDisplayName,
			},
			buildStubs: func(repo *mocks.MockUserRepository) {
				updatedUser := entity.User{
					DisplayName: newDisplayName,
				}
				user.DisplayName = newDisplayName

				repo.EXPECT().
					UpdateByID(gomock.Eq(user.ID), &updatedUser).
					Times(1).
					Return(&user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.UserDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, newDisplayName, response.Result.DisplayName)
			},
		},
		{
			name:       "Unauthorized",
			userID:     uuid.New().String(),
			body:       gin.H{},
			buildStubs: func(repo *mocks.MockUserRepository) {},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusUnauthorized), response.Status)
				require.Equal(t, http.StatusUnauthorized, response.StatusCode)
			},
		},
		{
			name:   "Internal Server Error",
			userID: user.ID.String(),
			body: gin.H{
				"display_name": newDisplayName,
			},
			buildStubs: func(repo *mocks.MockUserRepository) {
				updatedUser := entity.User{
					DisplayName: newDisplayName,
				}
				repo.EXPECT().
					UpdateByID(gomock.Eq(user.ID), &updatedUser).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
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
			tc.buildStubs(s.userRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.PATCH("users/:id", middleware.AuthMiddleware(s.tokenMaker), s.handler.UpdateUser)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range tc.body {
				err := writer.WriteField(key, value.(string))
				require.NoError(t, err)
			}

			err := writer.Close()
			require.NoError(t, err)

			url := fmt.Sprintf("/users/%s", tc.userID)
			request, err := http.NewRequest(http.MethodPatch, url, body)
			require.NoError(t, err)

			request.Header.Add("Content-Type", writer.FormDataContentType())
			c.Request = request

			tc.setupAuth(t, c.Request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
