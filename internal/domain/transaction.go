package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         string
	UserID     string
	CategoryID string
	Amount     decimal.Decimal
	Name       string
	RuleID     *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
