package entity

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	RefCode string `gorm:"size:255;not null;unique"`
}

type TransactionCreatePayload struct {
	RefCode string `json:"ref_code" binding:"required"`
} // @name TransactionCreatePayload

type TransactionDto struct {
	ID        uint   `json:"id"`
	RefCode   string `json:"ref_code"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
} // @name Transaction

func (t *Transaction) ToTransactionDto() *TransactionDto {
	return &TransactionDto{
		ID:        t.ID,
		RefCode:   t.RefCode,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}
