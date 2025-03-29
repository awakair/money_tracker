package dto_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"

	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
)

func TestFromTransactionRule(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domain.TransactionRule
		expected *dto.TransactionRuleDTO
	}{
		{
			name: "should correctly convert transaction rule without last generated time",
			input: &domain.TransactionRule{
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
			expected: &dto.TransactionRuleDTO{
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
			name: "should correctly convert transaction rule with last generated time",
			input: &domain.TransactionRule{
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
			expected: &dto.TransactionRuleDTO{
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
			result := dto.FromTransactionRule(tc.input)

			require.NotNil(t, result)
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
