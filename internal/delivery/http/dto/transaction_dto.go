package dto

import (
	"time"

	"github.com/shopspring/decimal"

	usecaseDTO "money_tracker/internal/application/usecases/dto"
)

// CreateTransactionRequest represents the data for creating a transaction
type CreateTransactionRequest struct {
	CategoryID string          `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount     decimal.Decimal `json:"amount" example:"99.99"`
	Name       string          `json:"name" example:"Weekly grocery shopping"`
}

func (r *CreateTransactionRequest) ToUsecaseDTO() usecaseDTO.CreateTransactionDTO {
	return usecaseDTO.CreateTransactionDTO{
		CategoryID: r.CategoryID,
		Amount:     r.Amount,
		Name:       r.Name,
	}
}

// UpdateTransactionRequest represents the data for updating a transaction
type UpdateTransactionRequest struct {
	CategoryID string          `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount     decimal.Decimal `json:"amount" example:"149.99"`
	Name       string          `json:"name" example:"Monthly grocery shopping"`
}

func (r *UpdateTransactionRequest) ToUsecaseDTO() usecaseDTO.UpdateTransactionDTO {
	return usecaseDTO.UpdateTransactionDTO{
		CategoryID: r.CategoryID,
		Amount:     r.Amount,
		Name:       r.Name,
	}
}

// TransactionResponse represents transaction information
type TransactionResponse struct {
	ID         string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     string          `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	CategoryID string          `json:"category_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Amount     decimal.Decimal `json:"amount" example:"99.99"`
	Name       string          `json:"name" example:"Weekly grocery shopping"`
	RuleID     *string         `json:"rule_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440003"`
	CreatedAt  time.Time       `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt  time.Time       `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func NewTransactionResponseFromUsecaseDTO(dto *usecaseDTO.TransactionDTO) TransactionResponse {
	return TransactionResponse{
		ID:         dto.ID,
		UserID:     dto.UserID,
		CategoryID: dto.CategoryID,
		Amount:     dto.Amount,
		Name:       dto.Name,
		RuleID:     dto.RuleID,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}

// TransactionListResponse represents a paginated list of transactions
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total" example:"50"`
	Limit        *int                  `json:"limit,omitempty" example:"20"`
	Offset       *int                  `json:"offset,omitempty" example:"0"`
}
