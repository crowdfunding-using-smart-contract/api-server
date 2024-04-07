package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HelloSuite struct {
	suite.Suite
}

func (s *HelloSuite) SetupTest() {

}

func (s *HelloSuite) TestGetHelloAPI() {
	type Query struct {
		Name string
	}

	testCases := []struct {
		name          string
		query         Query
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: Query{},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response MessageResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusOK, response.StatusCode)
				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, "Hello, Guest!", response.Message)
			},
		},
		{
			name: "OK with given name",
			query: Query{
				Name: "John",
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response MessageResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusOK, response.StatusCode)
				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, "Hello, John!", response.Message)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			url := "/hello"
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			q := c.Request.URL.Query()
			q.Add("name", tc.query.Name)
			c.Request.URL.RawQuery = q.Encode()

			r.GET("/hello", GetHelloMessage)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestHelloSuite(t *testing.T) {
	suite.Run(t, new(HelloSuite))
}
