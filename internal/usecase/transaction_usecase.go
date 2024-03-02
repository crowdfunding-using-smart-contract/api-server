package usecase

import (
	"fund-o/api-server/internal/datasource/repository"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/pkg/pagination"
)

type TransactionUseCase interface {
	CreateTransaction(transaction *entity.TransactionCreatePayload) (*entity.TransactionDto, error)
	GetTransactionByRefCode(refCode string) (*entity.TransactionDto, error)
	ListTransactions(paginateOptions pagination.PaginateOptions) pagination.PaginateResult[entity.TransactionDto]
}

type transactionRepository struct {
	transactionRepository repository.TransactionRepository
}

type TransactionUseCaseOptions struct {
	TransactionRepository repository.TransactionRepository
}

func NewTransactionUseCase(options *TransactionUseCaseOptions) TransactionUseCase {
	return &transactionRepository{
		transactionRepository: options.TransactionRepository,
	}
}

func (uc *transactionRepository) CreateTransaction(transaction *entity.TransactionCreatePayload) (*entity.TransactionDto, error) {
	newTransaction := &entity.Transaction{
		RefCode: transaction.RefCode,
	}
	newTransaction, err := uc.transactionRepository.Create(newTransaction)
	if err != nil {
		return nil, err
	}

	return newTransaction.ToTransactionDto(), nil
}

func (uc *transactionRepository) GetTransactionByRefCode(refCode string) (*entity.TransactionDto, error) {
	existingTransaction, err := uc.transactionRepository.FindByRefCode(refCode)
	if err != nil {
		return nil, err
	}
	return existingTransaction.ToTransactionDto(), nil
}

func (uc *transactionRepository) ListTransactions(paginateOptions pagination.PaginateOptions) pagination.PaginateResult[entity.TransactionDto] {
	result := pagination.MakePaginateResult(pagination.MakePaginateContextParameters[entity.TransactionDto]{
		PaginateOptions: paginateOptions,
		CountDocuments: func() int64 {
			return uc.transactionRepository.Count()
		},
		FindDocuments: func(findOptions pagination.PaginateFindOptions) []entity.TransactionDto {
			documents := uc.transactionRepository.List(findOptions)

			transactionDtos := make([]entity.TransactionDto, 0, len(documents))
			for _, document := range documents {
				transactionDtos = append(transactionDtos, *document.ToTransactionDto())
			}

			return transactionDtos
		},
	})

	return result
}
