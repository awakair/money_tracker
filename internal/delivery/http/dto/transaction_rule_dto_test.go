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

func TestCreateTransactionRuleRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.CreateTransactionRuleRequest
		expected dto.CreateTransactionRuleDTO
	}{
		{
			name: "should correctly convert create transaction rule request",
			input: &httpDTO.CreateTransactionRuleRequest{
				CategoryID:     "cat-123",
				Amount:         decimal.NewFromFloat(1000.00),
				Name:           "Monthly rent",
				CronExpression: "0 0 1 * *",
			},
			expected: dto.CreateTransactionRuleDTO{
				CategoryID:     "cat-123",
				Amount:         decimal.NewFromFloat(1000.00),
				Name:           "Monthly rent",
				CronExpression: "0 0 1 * *",
			},
		},
		{
			name: "should correctly convert create transaction rule request with empty fields",
			input: &httpDTO.CreateTransactionRuleRequest{
				CategoryID:     "",
				Amount:         decimal.Zero,
				Name:           "",
				CronExpression: "",
			},
			expected: dto.CreateTransactionRuleDTO{
				CategoryID:     "",
				Amount:         decimal.Zero,
				Name:           "",
				CronExpression: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.CronExpression, result.CronExpression)
		})
	}
}

func TestUpdateTransactionRuleRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.UpdateTransactionRuleRequest
		expected dto.UpdateTransactionRuleDTO
	}{
		{
			name: "should correctly convert update transaction rule request",
			input: &httpDTO.UpdateTransactionRuleRequest{
				CategoryID:     "cat-456",
				Amount:         decimal.NewFromFloat(1200.00),
				Name:           "Updated rent",
				CronExpression: "0 0 2 * *",
			},
			expected: dto.UpdateTransactionRuleDTO{
				CategoryID:     "cat-456",
				Amount:         decimal.NewFromFloat(1200.00),
				Name:           "Updated rent",
				CronExpression: "0 0 2 * *",
			},
		},
		{
			name: "should correctly convert update transaction rule request with different cron",
			input: &httpDTO.UpdateTransactionRuleRequest{
				CategoryID:     "cat-789",
				Amount:         decimal.NewFromFloat(50.00),
				Name:           "Weekly groceries",
				CronExpression: "0 0 * * 0",
			},
			expected: dto.UpdateTransactionRuleDTO{
				CategoryID:     "cat-789",
				Amount:         decimal.NewFromFloat(50.00),
				Name:           "Weekly groceries",
				CronExpression: "0 0 * * 0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.CronExpression, result.CronExpression)
		})
	}
}

func TestNewTransactionRuleResponseFromUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *dto.TransactionRuleDTO
		expected httpDTO.TransactionRuleResponse
	}{
		{
			name: "should correctly convert transaction rule DTO without last generated time",
			input: &dto.TransactionRuleDTO{
				ID:              "rule-123",
				UserID:          "user-456",
				CategoryID:      "cat-789",
				Amount:          decimal.NewFromFloat(1000.00),
				Name:            "Monthly rent",
				CronExpression:  "0 0 1 * *",
				LastGeneratedAt: nil,
				CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: httpDTO.TransactionRuleResponse{
				ID:              "rule-123",
				UserID:          "user-456",
				CategoryID:      "cat-789",
				Amount:          decimal.NewFromFloat(1000.00),
				Name:            "Monthly rent",
				CronExpression:  "0 0 1 * *",
				LastGeneratedAt: nil,
				CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "should correctly convert transaction rule DTO with last generated time",
			input: &dto.TransactionRuleDTO{
				ID:              "rule-456",
				UserID:          "user-789",
				CategoryID:      "cat-111",
				Amount:          decimal.NewFromFloat(50.00),
				Name:            "Weekly groceries",
				CronExpression:  "0 0 * * 0",
				LastGeneratedAt: ptr.To(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)),
				CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: httpDTO.TransactionRuleResponse{
				ID:              "rule-456",
				UserID:          "user-789",
				CategoryID:      "cat-111",
				Amount:          decimal.NewFromFloat(50.00),
				Name:            "Weekly groceries",
				CronExpression:  "0 0 * * 0",
				LastGeneratedAt: ptr.To(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)),
				CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := httpDTO.NewTransactionRuleResponseFromUsecaseDTO(tc.input)

			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.UserID, result.UserID)
			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.CronExpression, result.CronExpression)
			assert.Equal(t, tc.expected.LastGeneratedAt, result.LastGeneratedAt)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
