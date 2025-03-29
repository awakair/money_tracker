package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionRule struct {
	ID              string
	UserID          string
	CategoryID      string
	Amount          decimal.Decimal
	Name            string
	CronExpression  string
	LastGeneratedAt *time.Time
	DisabledAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (r *TransactionRule) ToTransaction(id string, now time.Time) *Transaction {
	return &Transaction{
		ID:         id,
		UserID:     r.UserID,
		CategoryID: r.CategoryID,
		Amount:     r.Amount,
		Name:       r.Name,
		RuleID:     &r.ID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
