package dto_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	"money_tracker/internal/application/usecases/dto"
	httpDTO "money_tracker/internal/delivery/http/dto"
)

func TestCreateTransactionRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.CreateTransactionRequest
		expected dto.CreateTransactionDTO
	}{
		{
			name: "should correctly convert create transaction request",
			input: &httpDTO.CreateTransactionRequest{
				CategoryID: "cat-123",
				Amount:     decimal.NewFromFloat(100.50),
				Name:       "Groceries",
			},
			expected: dto.CreateTransactionDTO{
				CategoryID: "cat-123",
				Amount:     decimal.NewFromFloat(100.50),
				Name:       "Groceries",
			},
		},
		{
			name: "should correctly convert create transaction request with zero amount",
			input: &httpDTO.CreateTransactionRequest{
				CategoryID: "cat-456",
				Amount:     decimal.Zero,
				Name:       "Test",
			},
			expected: dto.CreateTransactionDTO{
				CategoryID: "cat-456",
				Amount:     decimal.Zero,
				Name:       "Test",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
		})
	}
}

func TestUpdateTransactionRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.UpdateTransactionRequest
		expected dto.UpdateTransactionDTO
	}{
		{
			name: "should correctly convert update transaction request",
			input: &httpDTO.UpdateTransactionRequest{
				CategoryID: "cat-789",
				Amount:     decimal.NewFromFloat(200.75),
				Name:       "Updated Groceries",
			},
			expected: dto.UpdateTransactionDTO{
				CategoryID: "cat-789",
				Amount:     decimal.NewFromFloat(200.75),
				Name:       "Updated Groceries",
			},
		},
		{
			name: "should correctly convert update transaction request with empty name",
			input: &httpDTO.UpdateTransactionRequest{
				CategoryID: "cat-111",
				Amount:     decimal.NewFromFloat(50.00),
				Name:       "",
			},
			expected: dto.UpdateTransactionDTO{
				CategoryID: "cat-111",
				Amount:     decimal.NewFromFloat(50.00),
				Name:       "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
		})
	}
}

func TestNewTransactionResponseFromUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *dto.TransactionDTO
		expected httpDTO.TransactionResponse
	}{
		{
			name: "should correctly convert transaction DTO without rule",
			input: &dto.TransactionDTO{
				ID:         "txn-123",
				UserID:     "user-456",
				CategoryID: "cat-789",
				Amount:     decimal.NewFromFloat(100.50),
				Name:       "Groceries",
				RuleID:     nil,
				CreatedAt:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			expected: httpDTO.TransactionResponse{
				ID:         "txn-123",
				UserID:     "user-456",
				CategoryID: "cat-789",
				Amount:     decimal.NewFromFloat(100.50),
				Name:       "Groceries",
				RuleID:     nil,
				CreatedAt:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "should correctly convert transaction DTO with rule",
			input: &dto.TransactionDTO{
				ID:         "txn-456",
				UserID:     "user-789",
				CategoryID: "cat-111",
				Amount:     decimal.NewFromFloat(1500.00),
				Name:       "Monthly rent",
				RuleID:     ptr.To("rule-123"),
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: httpDTO.TransactionResponse{
				ID:         "txn-456",
				UserID:     "user-789",
				CategoryID: "cat-111",
				Amount:     decimal.NewFromFloat(1500.00),
				Name:       "Monthly rent",
				RuleID:     ptr.To("rule-123"),
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := httpDTO.NewTransactionResponseFromUsecaseDTO(tc.input)

			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.UserID, result.UserID)
			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.RuleID, result.RuleID)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
