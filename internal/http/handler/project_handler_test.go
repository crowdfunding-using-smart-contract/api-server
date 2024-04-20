package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/internal/usecase"
	"fund-o/api-server/mocks"
	"fund-o/api-server/pkg/pagination"
	"fund-o/api-server/pkg/random"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type ProjectTestSuite struct {
	suite.Suite
	tokenMaker                token.Maker
	repository                *mocks.MockProjectRepository
	projectCategoryRepository *mocks.MockProjectCategoryRepository
	handler                   *ProjectHandler
}

func (s *ProjectTestSuite) SetupSuite() {
	var err error

	s.tokenMaker, err = token.NewJWTMaker(random.NewString(32))
	require.NoError(s.T(), err)

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.repository = mocks.NewMockProjectRepository(ctrl)
	s.projectCategoryRepository = mocks.NewMockProjectCategoryRepository(ctrl)
	projectUseCase := usecase.NewProjectUseCase(&usecase.ProjectUseCaseOptions{
		ProjectRepository: s.repository,
	})
	projectCategoryUseCase := usecase.NewProjectCategoryUseCase(&usecase.ProjectCategoryUseCaseOptions{
		ProjectCategoryRepository: s.projectCategoryRepository,
	})
	s.handler = NewProjectHandler(&ProjectHandlerOptions{
		ProjectUseCase:         projectUseCase,
		ProjectCategoryUseCase: projectCategoryUseCase,
	})
}

func (s *ProjectTestSuite) TestListProjectsAPI() {
	projects := randomProjects(20)

	testCases := []struct {
		name          string
		query         entity.ProjectListParams
		buildStubs    func(repo *mocks.MockProjectRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: entity.ProjectListParams{
				PaginateOptions: pagination.PaginateOptions{
					Page: 1,
					Size: 10,
				},
				Query: "title",
			},
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					Count().
					Times(1).
					Return(int64(len(projects)))
				repo.EXPECT().
					FindAll(gomock.Any(), gomock.Any()).
					Times(1).
					Return(projects)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[pagination.PaginateResult[entity.PostDto]]

				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)
				require.NotNil(t, response.Result)

				require.Equal(t, 1, response.Result.CurrentPage)
				require.Equal(t, 10, response.Result.PerPage)
			},
		},
		{
			name: "Bad Request",
			query: entity.ProjectListParams{
				PaginateOptions: pagination.PaginateOptions{
					Page: -1,
				},
			},
			buildStubs: func(repo *mocks.MockProjectRepository) {},
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
			tc.buildStubs(s.repository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/projects", s.handler.ListProjects)

			url := "/projects"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page", strconv.Itoa(tc.query.Page))
			q.Add("size", strconv.Itoa(tc.query.Size))
			q.Add("q", tc.query.Query)

			c.Request = request
			c.Request.URL.RawQuery = q.Encode()

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ProjectTestSuite) TestGetProjectAPI() {
	project := randomProjects(1)[0]

	testCases := []struct {
		name          string
		projectID     string
		buildStubs    func(repo *mocks.MockProjectRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			projectID: uuid.NewString(),
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindByID(gomock.Any()).
					Times(1).
					Return(&project, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[entity.ProjectDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Equal(t, project.ID.String(), response.Result.ID)
			},
		},
		{
			name:       "Bad Request",
			projectID:  "invalid-uuid",
			buildStubs: func(repo *mocks.MockProjectRepository) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ErrorResponse

				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
				require.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
		},
		{
			name:      "Not Found",
			projectID: uuid.NewString(),
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindByID(gomock.Any()).
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
		{
			name:      "Internal Server Error",
			projectID: uuid.NewString(),
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindByID(gomock.Any()).
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
			tc.buildStubs(s.repository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/projects/:id", s.handler.GetProjectByID)

			url := fmt.Sprintf("/projects/%s", tc.projectID)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ProjectTestSuite) TestGetOwnProjectsAPI() {
	projects := randomProjects(10)
	user := randomUser(s.T())

	testCases := []struct {
		name          string
		userID        string
		buildStubs    func(repo *mocks.MockProjectRepository)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: uuid.NewString(),
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindAllByOwnerID(gomock.Any()).
					Times(1).
					Return(projects, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, s.tokenMaker, middleware.AuthorizationTypeBearer, user.ID.String(), time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[[]entity.ProjectDto]

				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Len(t, response.Result, len(projects))
			},
		},
		{
			name:   "Internal Server Error",
			userID: uuid.NewString(),
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindAllByOwnerID(gomock.Any()).
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
			tc.buildStubs(s.repository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/projects/me", middleware.AuthMiddleware(s.tokenMaker), s.handler.GetOwnProjects)

			url := "/projects/me"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			c.Request = request

			tc.setupAuth(t, c.Request, s.tokenMaker)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ProjectTestSuite) TestGetRecommendProjectsAPI() {
	projects := randomProjects(3)

	testCases := []struct {
		name          string
		buildStubs    func(repo *mocks.MockProjectRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindRecommendation(gomock.Any()).
					Times(1).
					Return(projects, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[[]entity.ProjectDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Len(t, response.Result, len(projects))
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(repo *mocks.MockProjectRepository) {
				repo.EXPECT().
					FindRecommendation(gomock.Any()).
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
			tc.buildStubs(s.repository)

			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			r.GET("/projects/recommendation", s.handler.GetRecommendProjects)

			url := "/projects/recommendation"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			c.Request = request

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *ProjectTestSuite) TestListProjectCategoriesAPI() {
	pc := []entity.ProjectCategory{
		{
			Name: "Category 1",
			SubCategories: []entity.ProjectSubCategory{
				{
					Name: "Sub Category 1",
				},
				{
					Name: "Sub Category 2",
				},
			},
		},
		{
			Name: "Category 2",
			SubCategories: []entity.ProjectSubCategory{
				{
					Name: "Sub Category 3",
				},
				{
					Name: "Sub Category 4",
				},
			},
		},
	}
	testCases := []struct {
		name          string
		buildStubs    func(repo *mocks.MockProjectCategoryRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mocks.MockProjectCategoryRepository) {
				repo.EXPECT().
					FindAll().
					Times(1).
					Return(pc, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var response ResultResponse[[]entity.ProjectCategoryDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusText(http.StatusOK), response.Status)
				require.Equal(t, http.StatusOK, response.StatusCode)

				require.NotNil(t, response.Result)
				require.Len(t, response.Result, len(pc))
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(repo *mocks.MockProjectCategoryRepository) {
				repo.EXPECT().
					FindAll().
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
			recorder := httptest.NewRecorder()
			c, r := gin.CreateTestContext(recorder)

			tc.buildStubs(s.projectCategoryRepository)

			r.GET("/projects/categories", s.handler.ListProjectCategories)

			url := "/projects/categories"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			c.Request = request

			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomProjects(n int) []entity.Project {
	projects := make([]entity.Project, n)
	for i := 0; i < n; i++ {
		projectRatings := make([]entity.ProjectRating, n)
		for j := 0; j < n; j++ {
			projectRatings[j] = entity.ProjectRating{
				Rating: random.NewFloat32(1, 5),
				UserID: uuid.New(),
			}
		}

		projects[i] = entity.Project{
			Base: entity.Base{
				ID: uuid.New(),
			},
			ProjectContractID: strconv.Itoa(random.NewInt(1, 1000)),
			Title:             fmt.Sprintf("title %d", i+1),
			SubTitle:          fmt.Sprintf("sub title %d", i+1),
			CategoryID:        uuid.New(),
			SubCategoryID:     uuid.New(),
			Location:          "Bangkok, Thailand",
			Ratings:           projectRatings,
			StartDate:         time.Now(),
			EndDate:           time.Now().AddDate(0, 0, 30),
			OwnerID:           uuid.New(),
		}
	}
	return projects
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
