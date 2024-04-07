package middleware

import (
	"fmt"
	"fund-o/api-server/pkg/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MiddlewareSuite struct {
	suite.Suite
	tokenMaker token.Maker
}

func (s *MiddlewareSuite) SetupTest() {
	var err error
	secretKey := "alsypVB6YUpE2HBW4npGoXeArNyqVrqO"

	s.tokenMaker, err = token.NewJWTMaker(secretKey)
	require.NoError(s.T(), err)
}

func (s *MiddlewareSuite) TestAuthorizationMiddleware() {
	userID := uuid.New().String()

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				s.addAuthorization(t, request, AuthorizationTypeBearer, userID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				s.addAuthorization(t, request, "unsupported", userID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				s.addAuthorization(t, request, "", userID, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				s.addAuthorization(t, request, AuthorizationTypeBearer, userID, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/auth",
				AuthMiddleware(s.tokenMaker),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{})
				},
			)

			request, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			c.Request = request

			tc.setupAuth(t, request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)

			tc.checkResponse(t, recorder)
		})
	}
}

func (s *MiddlewareSuite) addAuthorization(
	t *testing.T,
	request *http.Request,
	authorizationType string,
	userID string,
	duration time.Duration,
) {
	token, payload, err := s.tokenMaker.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(AuthorizationHeaderKey, authorizationHeader)
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
