package dto_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"

	"k8s.io/utils/ptr"
)

func TestFromTransaction(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domain.Transaction
		expected *dto.TransactionDTO
	}{
		{
			name: "should correctly convert transaction without rule",
			input: &domain.Transaction{
				ID:         "txn-123",
				UserID:     "user-456",
				CategoryID: "cat-789",
				Amount:     decimal.NewFromFloat(100.50),
				Name:       "Groceries",
				RuleID:     nil,
				CreatedAt:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			expected: &dto.TransactionDTO{
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
			name: "should correctly convert transaction with rule",
			input: &domain.Transaction{
				ID:         "txn-456",
				UserID:     "user-789",
				CategoryID: "cat-111",
				Amount:     decimal.NewFromFloat(1500.00),
				Name:       "Monthly rent",
				RuleID:     ptr.To("rule-123"),
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: &dto.TransactionDTO{
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
			result := dto.FromTransaction(tc.input)

			require.NotNil(t, result)
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
