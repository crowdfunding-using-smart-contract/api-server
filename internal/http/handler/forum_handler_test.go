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
	"fund-o/api-server/pkg/pagination"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type ForumTestSuite struct {
	suite.Suite
	tokenMaker      token.Maker
	forumRepository *mocks.MockForumRepository
	handler         *ForumHandler
}

func (s *ForumTestSuite) SetupSuite() {
	var err error
	secretKey := "alsypVB6YUpE2HBW4npGoXeArNyqVrqO"

	s.tokenMaker, err = token.NewJWTMaker(secretKey)
	s.Require().NoError(err)

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.forumRepository = mocks.NewMockForumRepository(ctrl)
	forumUseCase := usecase.NewForumUseCase(&usecase.ForumUseCaseOptions{
		ForumRepository: s.forumRepository,
	})
	s.handler = NewForumHandler(&ForumHandlerOptions{
		ForumUseCase: forumUseCase,
	})
}

func (s *ForumTestSuite) TestListPostsAPI() {
	posts := randomPosts(35)

	testCases := []struct {
		name          string
		query         pagination.PaginateOptions
		buildStubs    func(repo *mocks.MockForumRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: pagination.PaginateOptions{},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CountPost().
					Times(1).
					Return(int64(len(posts)))
				repo.EXPECT().
					ListPosts(gomock.Any()).
					Times(1).
					Return(posts)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[pagination.PaginateResult[entity.PostDto]]

				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, 1, response.Result.CurrentPage)
				require.Equal(t, 20, response.Result.PerPage)
			},
		},
		{
			name: "OK with custom pagination",
			query: pagination.PaginateOptions{
				Page: 2,
				Size: 10,
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CountPost().
					Times(1).
					Return(int64(len(posts)))
				repo.EXPECT().
					ListPosts(gomock.Any()).
					Times(1).
					Return(posts[10:20])
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[pagination.PaginateResult[entity.PostDto]]

				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, 2, response.Result.CurrentPage)
				require.Equal(t, 10, response.Result.PerPage)
			},
		},
		{
			name:  "OK with empty posts",
			query: pagination.PaginateOptions{},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CountPost().
					Times(1).
					Return(int64(0))
				repo.EXPECT().
					ListPosts(gomock.Any()).
					Times(1).
					Return([]entity.Post{})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[pagination.PaginateResult[entity.PostDto]]

				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, 1, response.Result.CurrentPage)
				require.Equal(t, 20, response.Result.PerPage)
				require.Equal(t, int64(0), response.Result.Total)
				require.Equal(t, 0, len(response.Result.Data))
			},
		},
		{
			name: "Bad Request",
			query: pagination.PaginateOptions{
				Page: -1,
			},
			buildStubs: func(repo *mocks.MockForumRepository) {},
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
			tc.buildStubs(s.forumRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/posts", s.handler.ListPosts)

			url := "/posts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page", strconv.Itoa(tc.query.Page))
			q.Add("size", strconv.Itoa(tc.query.Size))

			c.Request = request
			c.Request.URL.RawQuery = q.Encode()

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ForumTestSuite) TestGetPostAPI() {
	post := randomPosts(1)[0]

	testCases := []struct {
		name          string
		postID        string
		buildStubs    func(repo *mocks.MockForumRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			postID: post.ID.String(),
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					FindPostByID(gomock.Eq(post.ID)).
					Times(1).
					Return(&post, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.PostDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Equal(t, post.ID.String(), response.Result.ID)
			},
		},
		{
			name:   "Internal Server Error",
			postID: post.ID.String(),
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					FindPostByID(gomock.Eq(post.ID)).
					Times(1).
					Return(nil, sql.ErrConnDone)
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
			tc.buildStubs(s.forumRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/posts/:id", s.handler.GetPostByID)

			url := fmt.Sprintf("/posts/%s", tc.postID)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ForumTestSuite) TestCreatePostAPI() {
	user := randomUser(s.T())
	projectID := uuid.New()

	testCases := []struct {
		name          string
		payload       gin.H
		buildStubs    func(repo *mocks.MockForumRepository)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			payload: gin.H{
				"title":       "Post 1",
				"description": "Description 1",
				"content":     "Content 1",
				"project_id":  projectID.String(),
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				post := entity.Post{
					Title:       "Post 1",
					Description: "Description 1",
					Content:     "Content 1",
					ProjectID:   projectID,
					AuthorID:    user.ID,
				}
				repo.EXPECT().
					CreatePost(&post).
					Times(1).
					Return(&post, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.PostDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusCreated), response.Status)
				require.Equal(t, http.StatusCreated, response.StatusCode)
			},
		},
		{
			name:    "Bad Request",
			payload: gin.H{},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CreatePost(gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
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
			name: "Internal Server Error",
			payload: gin.H{
				"title":       "Post 1",
				"description": "Description 1",
				"content":     "Content 1",
				"project_id":  projectID.String(),
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CreatePost(gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
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
			tc.buildStubs(s.forumRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.POST("/posts", middleware.AuthMiddleware(s.tokenMaker), s.handler.CreatePost)

			requestBody, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			url := "/posts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
			require.NoError(t, err)

			c.Request = request

			tc.setupAuth(t, c.Request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ForumTestSuite) TestCreateCommentAPI() {
	user := randomUser(s.T())

	testCases := []struct {
		name          string
		postID        string
		payload       gin.H
		buildStubs    func(repo *mocks.MockForumRepository)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			postID: uuid.NewString(),
			payload: gin.H{
				"content": "Hello, Bro!",
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				comment := entity.Comment{
					Content:  "Hello, Bro!",
					AuthorID: user.ID,
				}
				repo.EXPECT().
					CreateComment(gomock.Any()).
					Times(1).
					Return(&comment, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.CommentDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusCreated), response.Status)
				require.Equal(t, http.StatusCreated, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Equal(t, response.Result.Content, "Hello, Bro!")
				require.Len(t, response.Result.Replies, 0)
			},
		},
		{
			name:       "Invalid Request Body",
			postID:     uuid.NewString(),
			payload:    gin.H{},
			buildStubs: func(repo *mocks.MockForumRepository) {},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
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
			name:   "Internal Server Error",
			postID: uuid.NewString(),
			payload: gin.H{
				"content": "Hello, Bro!",
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CreateComment(gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
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
			tc.buildStubs(s.forumRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.POST("/posts/:id/comments", middleware.AuthMiddleware(s.tokenMaker), s.handler.CreateComment)

			requestBody, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			url := fmt.Sprintf("/posts/%s/comments", tc.postID)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
			require.NoError(t, err)

			c.Request = request

			tc.setupAuth(t, c.Request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ForumTestSuite) TestCreateReplyAPI() {
	user := randomUser(s.T())

	testCases := []struct {
		name          string
		commentID     string
		payload       gin.H
		buildStubs    func(repo *mocks.MockForumRepository)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			commentID: uuid.NewString(),
			payload: gin.H{
				"content": "How are you?",
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				reply := entity.Reply{
					Content:  "How are you?",
					AuthorID: user.ID,
				}
				repo.EXPECT().
					CreateReply(gomock.Any()).
					Times(1).
					Return(&reply, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.ReplyDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusCreated), response.Status)
				require.Equal(t, http.StatusCreated, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Equal(t, response.Result.Content, "How are you?")
			},
		},
		{
			name:       "Invalid Request Body",
			commentID:  uuid.NewString(),
			payload:    gin.H{},
			buildStubs: func(repo *mocks.MockForumRepository) {},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
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
			name:      "Internal Server Error",
			commentID: uuid.NewString(),
			payload: gin.H{
				"content": "How are you?",
			},
			buildStubs: func(repo *mocks.MockForumRepository) {
				repo.EXPECT().
					CreateReply(gomock.Any()).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), 5*time.Minute)
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
			tc.buildStubs(s.forumRepository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.POST("/comments/:id/replies", middleware.AuthMiddleware(s.tokenMaker), s.handler.CreateReply)

			requestBody, err := json.Marshal(tc.payload)
			require.NoError(t, err)

			url := fmt.Sprintf("/comments/%s/replies", tc.commentID)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
			require.NoError(t, err)

			c.Request = request

			tc.setupAuth(t, c.Request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomPosts(n int) []entity.Post {
	posts := make([]entity.Post, n)
	for i := range posts {
		posts[i] = entity.Post{
			Base: entity.Base{
				ID: uuid.New(),
			},
			Title:       fmt.Sprintf("Post %d", i+1),
			Description: fmt.Sprintf("Content %d", i+1),
			Content:     fmt.Sprintf("Content %d", i+1),
			AuthorID:    uuid.New(),
			ProjectID:   uuid.New(),
		}
	}
	return posts
}

func TestForumSuite(t *testing.T) {
	suite.Run(t, new(ForumTestSuite))
}
