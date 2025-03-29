package dto

import (
	"time"

	"github.com/shopspring/decimal"

	usecaseDTO "money_tracker/internal/application/usecases/dto"
)

// CreateTransactionRuleRequest represents the data for creating a transaction rule
type CreateTransactionRuleRequest struct {
	CategoryID     string          `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount         decimal.Decimal `json:"amount" example:"1500.00"`
	Name           string          `json:"name" example:"Monthly salary"`
	CronExpression string          `json:"cron_expression" example:"0 0 1 * *"`
}

func (r *CreateTransactionRuleRequest) ToUsecaseDTO() usecaseDTO.CreateTransactionRuleDTO {
	return usecaseDTO.CreateTransactionRuleDTO{
		CategoryID:     r.CategoryID,
		Amount:         r.Amount,
		Name:           r.Name,
		CronExpression: r.CronExpression,
	}
}

// UpdateTransactionRuleRequest represents the data for updating a transaction rule
type UpdateTransactionRuleRequest struct {
	CategoryID     string          `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount         decimal.Decimal `json:"amount" example:"1600.00"`
	Name           string          `json:"name" example:"Monthly salary (updated)"`
	CronExpression string          `json:"cron_expression" example:"0 0 1 * *"`
}

func (r *UpdateTransactionRuleRequest) ToUsecaseDTO() usecaseDTO.UpdateTransactionRuleDTO {
	return usecaseDTO.UpdateTransactionRuleDTO{
		CategoryID:     r.CategoryID,
		Amount:         r.Amount,
		Name:           r.Name,
		CronExpression: r.CronExpression,
	}
}

// TransactionRuleResponse represents transaction rule information
type TransactionRuleResponse struct {
	ID              string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID          string          `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	CategoryID      string          `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Amount          decimal.Decimal `json:"amount" example:"1500.00"`
	Name            string          `json:"name" example:"Monthly salary"`
	CronExpression  string          `json:"cron_expression" example:"0 0 1 * *"`
	LastGeneratedAt *time.Time      `json:"last_generated_at,omitempty" example:"2024-01-01T00:00:00Z"`
	CreatedAt       time.Time       `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt       time.Time       `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func NewTransactionRuleResponseFromUsecaseDTO(dto *usecaseDTO.TransactionRuleDTO) TransactionRuleResponse {
	return TransactionRuleResponse{
		ID:              dto.ID,
		UserID:          dto.UserID,
		CategoryID:      dto.CategoryID,
		Amount:          dto.Amount,
		Name:            dto.Name,
		CronExpression:  dto.CronExpression,
		LastGeneratedAt: dto.LastGeneratedAt,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

// TransactionRuleListResponse represents a paginated list of transaction rules
type TransactionRuleListResponse struct {
	Rules  []TransactionRuleResponse `json:"rules"`
	Total  int                       `json:"total" example:"5"`
	Limit  *int                      `json:"limit,omitempty" example:"20"`
	Offset *int                      `json:"offset,omitempty" example:"0"`
}
