package repository

import (
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/pagination"
	"github.com/rs/zerolog"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) (*entity.Transaction, error)
	FindByRefCode(refCode string) (*entity.Transaction, error)
	List(findOptions pagination.PaginateFindOptions) []entity.Transaction
	Count() int64
}

type transactionRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	logger := log.With().Str("module", "transaction_repository").Logger()
	return &transactionRepository{db, logger}
}

func (repo *transactionRepository) Create(transaction *entity.Transaction) (*entity.Transaction, error) {
	if result := repo.db.Create(&transaction); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to create transaction")
		return nil, result.Error
	}

	return transaction, nil
}

func (repo *transactionRepository) FindByRefCode(refCode string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if result := repo.db.Where("ref_code = ?", refCode).First(&transaction); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to find transaction by reference code: " + refCode)
		return nil, result.Error
	}

	return &transaction, nil
}

func (repo *transactionRepository) List(findOptions pagination.PaginateFindOptions) (transactions []entity.Transaction) {
	if result := repo.db.Limit(findOptions.Limit).Offset(findOptions.Skip).Find(&transactions); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to list transactions")
		return
	}

	return transactions
}

func (repo *transactionRepository) Count() int64 {
	var count int64
	if result := repo.db.Model(&entity.Transaction{}).Count(&count); result.Error != nil {
		repo.logger.Error().Err(result.Error).Msg("failed to count transactions")
		return 0
	}

	return count
}
