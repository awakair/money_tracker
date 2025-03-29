package dto

import (
	"time"

	"github.com/shopspring/decimal"

	"money_tracker/internal/domain"
)

type CreateTransactionDTO struct {
	CategoryID string
	Amount     decimal.Decimal
	Name       string
}

type UpdateTransactionDTO struct {
	CategoryID string
	Amount     decimal.Decimal
	Name       string
}

type TransactionDTO struct {
	ID         string
	UserID     string
	CategoryID string
	Amount     decimal.Decimal
	Name       string
	RuleID     *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func FromTransaction(transaction *domain.Transaction) *TransactionDTO {
	return &TransactionDTO{
		ID:         transaction.ID,
		UserID:     transaction.UserID,
		CategoryID: transaction.CategoryID,
		Amount:     transaction.Amount,
		Name:       transaction.Name,
		RuleID:     transaction.RuleID,
		CreatedAt:  transaction.CreatedAt,
		UpdatedAt:  transaction.UpdatedAt,
	}
}
