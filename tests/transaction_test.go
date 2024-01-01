package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockRepository "fund-o/api-server/internal/datasource/repository/mock"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/handler"
	"fund-o/api-server/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TransactionSuite struct {
	suite.Suite
	repository *mockRepository.MockTransactionRepository
	usecase    usecase.TransactionUsecase
	handler    *handler.TransactionHandler
}

func (s *TransactionSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.repository = mockRepository.NewMockTransactionRepository(ctrl)
	s.usecase = usecase.NewTransactionUsecase(&usecase.TransactionUsecaseOptions{
		TransactionRepository: s.repository,
	})
	s.handler = handler.NewTransactionHandler(&handler.TransactionHandlerOptions{
		TransactionUsecase: s.usecase,
	})
}

func (s *TransactionSuite) TestGetTransactionAPI() {
	transaction := entity.Transaction{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		RefCode: "azc-5301",
	}

	testCases := []struct {
		name          string
		refCode       string
		buildStubs    func(store *mockRepository.MockTransactionRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			refCode: transaction.RefCode,
			buildStubs: func(store *mockRepository.MockTransactionRepository) {
				store.EXPECT().
					FindByRefCode(gomock.Eq(transaction.RefCode)).
					Times(1).
					Return(&transaction, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response handler.ResultResponse[entity.TransactionDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.NotEmpty(t, response.Result)
				require.Equal(t, transaction.RefCode, response.Result.RefCode)
			},
		},
		{
			name:    "Not Found",
			refCode: transaction.RefCode,
			buildStubs: func(store *mockRepository.MockTransactionRepository) {
				store.EXPECT().
					FindByRefCode(gomock.Eq(transaction.RefCode)).
					Times(1).
					Return(nil, gorm.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

				var response handler.ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				require.NoError(t, err)

				require.Equal(t, http.StatusNotFound, response.StatusCode)
				require.Equal(t, http.StatusText(http.StatusNotFound), response.Status)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.buildStubs(s.repository)

			recorder := httptest.NewRecorder()
			c, r := NewTestContext(t, recorder)

			url := fmt.Sprintf("/transactions/%s", tc.refCode)
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			r.GET("/transactions/:id", s.handler.GetTransaction)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func (s *TransactionSuite) TestCreateTransactionAPI() {
	refCode := "azc-5301"

	testCases := []struct {
		name          string
		payload       entity.TransactionCreatePayload
		buildStubs    func(store *mockRepository.MockTransactionRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			payload: entity.TransactionCreatePayload{
				RefCode: refCode,
			},
			buildStubs: func(store *mockRepository.MockTransactionRepository) {
				transaction := entity.Transaction{
					RefCode: refCode,
				}

				store.EXPECT().
					Create(&transaction).
					Times(1).
					Return(&transaction, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				var response handler.ResultResponse[entity.TransactionDto]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.NotEmpty(t, response.Result)
				require.Equal(t, "azc-5301", response.Result.RefCode)
			},
		},
		{
			name:    "Bad Request",
			payload: entity.TransactionCreatePayload{},
			buildStubs: func(store *mockRepository.MockTransactionRepository) {
				store.EXPECT().
					Create(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

				var response handler.ErrorResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)

				require.Equal(t, http.StatusBadRequest, response.StatusCode)
				require.Equal(t, http.StatusText(http.StatusBadRequest), response.Status)
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.buildStubs(s.repository)

			recorder := httptest.NewRecorder()
			c, r := NewTestContext(t, recorder)

			url := "/transactions"
			payload, err := json.Marshal(tc.payload)
			require.NoError(t, err)
			reader := strings.NewReader(string(payload))
			require.NoError(t, err)
			c.Request = httptest.NewRequest(http.MethodPost, url, reader)

			r.POST("/transactions", s.handler.CreateTransaction)
			r.ServeHTTP(recorder, c.Request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}
