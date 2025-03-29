package dto

import (
	"time"

	"github.com/shopspring/decimal"

	"money_tracker/internal/domain"
)

type CreateTransactionRuleDTO struct {
	CategoryID     string
	Amount         decimal.Decimal
	Name           string
	CronExpression string
}

type UpdateTransactionRuleDTO struct {
	CategoryID     string
	Amount         decimal.Decimal
	Name           string
	CronExpression string
}

type TransactionRuleDTO struct {
	ID              string
	UserID          string
	CategoryID      string
	Amount          decimal.Decimal
	Name            string
	CronExpression  string
	LastGeneratedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func FromTransactionRule(rule *domain.TransactionRule) *TransactionRuleDTO {
	return &TransactionRuleDTO{
		ID:              rule.ID,
		UserID:          rule.UserID,
		CategoryID:      rule.CategoryID,
		Amount:          rule.Amount,
		Name:            rule.Name,
		CronExpression:  rule.CronExpression,
		LastGeneratedAt: rule.LastGeneratedAt,
		CreatedAt:       rule.CreatedAt,
		UpdatedAt:       rule.UpdatedAt,
	}
}
