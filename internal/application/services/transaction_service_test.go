package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/interfaces/mocks"
	"money_tracker/internal/application/services"
	servicemocks "money_tracker/internal/application/services/mocks"
	"money_tracker/internal/domain"
)

func TestTransactionService_GenerateTransactionFromRule(t *testing.T) {
	ctx := context.Background()
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	t.Run("should not generate transaction when next generation time is in future", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDB := mocks.NewMockDatabase(ctrl)
		mockCron := servicemocks.NewMockCron(ctrl)
		mockIDGen := mocks.NewMockIDGenerator(ctrl)
		clock := clockwork.NewFakeClockAt(fixedTime)

		service := services.NewTransactionService(mockDB, clock, mockCron, mockIDGen)

		rule := &domain.TransactionRule{
			ID:              "rule-1",
			UserID:          "user-1",
			CategoryID:      "cat-1",
			Amount:          decimal.NewFromFloat(100.50),
			Name:            "Monthly subscription",
			CronExpression:  "0 0 1 * *",
			LastGeneratedAt: nil,
			CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		futureTime := fixedTime.Add(24 * time.Hour)
		mockCron.EXPECT().
			GetNextGenerationTime("0 0 1 * *", rule.CreatedAt).
			Return(futureTime, nil)

		err := service.GenerateTransactionFromRule(ctx, rule)

		require.NoError(t, err)
	})

	t.Run("should generate transaction when next generation time has passed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDB := mocks.NewMockDatabase(ctrl)
		mockCron := servicemocks.NewMockCron(ctrl)
		mockIDGen := mocks.NewMockIDGenerator(ctrl)
		clock := clockwork.NewFakeClockAt(fixedTime)

		service := services.NewTransactionService(mockDB, clock, mockCron, mockIDGen)

		rule := &domain.TransactionRule{
			ID:              "rule-1",
			UserID:          "user-1",
			CategoryID:      "cat-1",
			Amount:          decimal.NewFromFloat(100.50),
			Name:            "Monthly subscription",
			CronExpression:  "0 0 1 * *",
			LastGeneratedAt: nil,
			CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		pastTime := fixedTime.Add(-24 * time.Hour)
		mockCron.EXPECT().
			GetNextGenerationTime("0 0 1 * *", rule.CreatedAt).
			Return(pastTime, nil)

		mockIDGen.EXPECT().
			Generate().
			Return("transaction-1")

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				mockTransactionRepo := mocks.NewMockTransactionRepositoryRW(ctrl)
				mockRuleRepo := mocks.NewMockTransactionRuleRepositoryRW(ctrl)

				mockTransactionRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, txn *domain.Transaction) error {
						assert.Equal(t, "transaction-1", txn.ID)
						assert.Equal(t, "user-1", txn.UserID)
						assert.Equal(t, "cat-1", txn.CategoryID)
						assert.Equal(t, decimal.NewFromFloat(100.50), txn.Amount)
						assert.Equal(t, "Monthly subscription", txn.Name)
						assert.NotNil(t, txn.RuleID)
						assert.Equal(t, "rule-1", *txn.RuleID)
						assert.Equal(t, fixedTime, txn.CreatedAt)
						assert.Equal(t, fixedTime, txn.UpdatedAt)
						return nil
					})

				mockRuleRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, r *domain.TransactionRule) error {
						assert.NotNil(t, r.LastGeneratedAt)
						assert.Equal(t, fixedTime, *r.LastGeneratedAt)
						return nil
					})

				repo := interfaces.RepositoriesReadWrite{
					TransactionRepositoryRW:     mockTransactionRepo,
					TransactionRuleRepositoryRW: mockRuleRepo,
				}
				return fn(ctx, repo)
			})

		err := service.GenerateTransactionFromRule(ctx, rule)

		require.NoError(t, err)
	})

	t.Run("should use LastGeneratedAt when available", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDB := mocks.NewMockDatabase(ctrl)
		mockCron := servicemocks.NewMockCron(ctrl)
		mockIDGen := mocks.NewMockIDGenerator(ctrl)
		clock := clockwork.NewFakeClockAt(fixedTime)

		service := services.NewTransactionService(mockDB, clock, mockCron, mockIDGen)

		lastGenTime := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
		rule := &domain.TransactionRule{
			ID:              "rule-1",
			UserID:          "user-1",
			CategoryID:      "cat-1",
			Amount:          decimal.NewFromFloat(100.50),
			Name:            "Monthly subscription",
			CronExpression:  "0 0 1 * *",
			LastGeneratedAt: &lastGenTime,
			CreatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		pastTime := fixedTime.Add(-24 * time.Hour)
		mockCron.EXPECT().
			GetNextGenerationTime("0 0 1 * *", lastGenTime).
			Return(pastTime, nil)

		mockIDGen.EXPECT().
			Generate().
			Return("transaction-1")

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				mockTransactionRepo := mocks.NewMockTransactionRepositoryRW(ctrl)
				mockRuleRepo := mocks.NewMockTransactionRuleRepositoryRW(ctrl)

				mockTransactionRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)

				mockRuleRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil)

				repo := interfaces.RepositoriesReadWrite{
					TransactionRepositoryRW:     mockTransactionRepo,
					TransactionRuleRepositoryRW: mockRuleRepo,
				}
				return fn(ctx, repo)
			})

		err := service.GenerateTransactionFromRule(ctx, rule)

		require.NoError(t, err)
	})
}
