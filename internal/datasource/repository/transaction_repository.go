package repository

import (
	"github.com/danyouknowme/gin-gorm-boilerplate/internal/entity"
	"github.com/danyouknowme/gin-gorm-boilerplate/pkg/pagination"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) (*entity.Transaction, error)
	FindByRefCode(refCode string) (*entity.Transaction, error)
	List(findOptions pagination.PaginateFindOptions) []entity.Transaction
	Count() int64
}

type transactionRepository struct {
	db *gorm.DB
}

var logger *log.Entry

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	logger = log.WithFields(log.Fields{
		"module": "transaction_repository",
	})
	return &transactionRepository{db}
}

func (repo *transactionRepository) Create(transaction *entity.Transaction) (*entity.Transaction, error) {
	if result := repo.db.Create(&transaction); result.Error != nil {
		logger.Errorf("Failed to create transaction: %v", result.Error)
		return nil, result.Error
	}

	return transaction, nil
}

func (repo *transactionRepository) FindByRefCode(refCode string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if result := repo.db.Where("ref_code = ?", refCode).First(&transaction); result.Error != nil {
		logger.Errorf("Failed to get transaction by reference code: %v", result.Error)
		return nil, result.Error
	}

	return &transaction, nil
}

func (repo *transactionRepository) List(findOptions pagination.PaginateFindOptions) (transactions []entity.Transaction) {
	// var transactions []entity.Transaction
	if result := repo.db.Limit(findOptions.Limit).Offset(findOptions.Skip).Find(&transactions); result.Error != nil {
		logger.Errorf("Failed to list transactions: %v", result.Error)
		return
	}

	return transactions
}

func (repo *transactionRepository) Count() int64 {
	var count int64
	if result := repo.db.Model(&entity.Transaction{}).Count(&count); result.Error != nil {
		logger.Errorf("Failed to count transactions: %v", result.Error)
		return 0
	}

	return count
}
