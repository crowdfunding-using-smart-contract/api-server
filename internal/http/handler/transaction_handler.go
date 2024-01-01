package handler

import (
	"fmt"
	"net/http"

	"github.com/danyouknowme/gin-gorm-boilerplate/internal/entity"
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/usecase"
	"github.com/danyouknowme/gin-gorm-boilerplate/pkg/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandlerOptions struct {
	TransactionUsecase usecase.TransactionUsecase
}

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(options *TransactionHandlerOptions) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: options.TransactionUsecase,
	}
}

// CreateTransaction godoc
// @summary Create Transaction
// @description Create transaction with reference code
// @tags transasctions
// @id CreateTransaction
// @accpet json
// @produce json
// @param Transaction body entity.TransactionCreatePayload true "Transaction data to be created"
// @response 200 {object} handler.ResultResponse[entity.TransactionDto] "OK"
// @response 400 {object} handler.ErrorResponse "Bad Request"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /transactions [post]
func (handler *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction entity.TransactionCreatePayload
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusBadRequest, fmt.Sprintf("error create transaction: %v", err.Error())))
		return
	}

	transactionDto, err := handler.transactionUsecase.CreateTransaction(&transaction)
	if err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error create transaction: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusCreated, transactionDto))
}

// GetTransaction godoc
// @summary Get Transaction
// @description Get transaction by id
// @tags transasctions
// @id GetTransaction
// @produce json
// @param id path string true "reference code of transaction"
// @response 200 {object} handler.ResultResponse[entity.TransactionDto] "OK"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /transactions/{id} [get]
func (handler *TransactionHandler) GetTransaction(c *gin.Context) {
	transactionRefCode := c.Param("id")

	transactionDto, err := handler.transactionUsecase.GetTransactionByRefCode(transactionRefCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(makeHttpErrorResponse(http.StatusNotFound, fmt.Sprintf("error get transaction: %v", err.Error())))
			return
		}

		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error get transaction: %v", err.Error())))
		return
	}

	c.JSON(makeHttpResponse(http.StatusOK, transactionDto))
}

// ListTransactions godoc
// @summary List Transaction
// @description Get list of transactions
// @tags transasctions
// @id ListTransactions
// @produce json
// @param page query int false "number of page"
// @param size query int false "size of data per page"
// @response 200 {object} handler.ResultResponse[pagination.PaginateResult[entity.TransactionDto]] "OK"
// @response 500 {object} handler.ErrorResponse "Internal Server Error"
// @router /transactions [get]
func (handler *TransactionHandler) ListTransactions(c *gin.Context) {
	var paginateOptions pagination.PaginateOptions
	if err := c.ShouldBindQuery(&paginateOptions); err != nil {
		c.JSON(makeHttpErrorResponse(http.StatusInternalServerError, fmt.Sprintf("error get transaction: %v", err.Error())))
		return
	}

	transactions := handler.transactionUsecase.ListTransactions(paginateOptions)
	c.JSON(makeHttpResponse(http.StatusOK, transactions))
}
