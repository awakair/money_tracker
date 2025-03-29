package domain_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money_tracker/internal/domain"

	"k8s.io/utils/ptr"
)

func TestTransactionRule_ToTransaction(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		rule     *domain.TransactionRule
		id       string
		now      time.Time
		expected *domain.Transaction
	}{
		{
			name: "should correctly convert transaction rule to transaction",
			rule: &domain.TransactionRule{
				ID:             "rule-123",
				UserID:         "user-456",
				CategoryID:     "cat-789",
				Amount:         decimal.NewFromFloat(1000.00),
				Name:           "Monthly rent",
				CronExpression: "0 0 1 * *",
			},
			id:  "txn-123",
			now: fixedTime,
			expected: &domain.Transaction{
				ID:         "txn-123",
				UserID:     "user-456",
				CategoryID: "cat-789",
				Amount:     decimal.NewFromFloat(1000.00),
				Name:       "Monthly rent",
				RuleID:     ptr.To("rule-123"),
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
		},
		{
			name: "should correctly convert transaction rule with different values",
			rule: &domain.TransactionRule{
				ID:             "rule-456",
				UserID:         "user-789",
				CategoryID:     "cat-111",
				Amount:         decimal.NewFromFloat(50.50),
				Name:           "Weekly groceries",
				CronExpression: "0 0 * * 0",
			},
			id:  "txn-456",
			now: time.Date(2025, 2, 1, 10, 30, 0, 0, time.UTC),
			expected: &domain.Transaction{
				ID:         "txn-456",
				UserID:     "user-789",
				CategoryID: "cat-111",
				Amount:     decimal.NewFromFloat(50.50),
				Name:       "Weekly groceries",
				RuleID:     ptr.To("rule-456"),
				CreatedAt:  time.Date(2025, 2, 1, 10, 30, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 2, 1, 10, 30, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.rule.ToTransaction(tc.id, tc.now)

			require.NotNil(t, result)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.UserID, result.UserID)
			assert.Equal(t, tc.expected.CategoryID, result.CategoryID)
			assert.True(t, tc.expected.Amount.Equal(result.Amount))
			assert.Equal(t, tc.expected.Name, result.Name)
			require.NotNil(t, result.RuleID)
			assert.Equal(t, *tc.expected.RuleID, *result.RuleID)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
