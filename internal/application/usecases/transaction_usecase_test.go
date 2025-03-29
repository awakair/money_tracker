package usecases_test

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
	"money_tracker/internal/application/usecases"
	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
	"money_tracker/internal/pkg/credscontext"
)

func TestTransactionUseCase_CreateTransaction(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepo := mocks.NewMockCategoryRepositoryRW(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionRepositoryRW(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionDTO{
			Name:       "Groceries",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
		}

		result, err := usecase.CreateTransaction(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionDTO{
			Name:       "   ",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
		}

		result, err := usecase.CreateTransaction(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when category_id is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionDTO{
			Name:       "Groceries",
			CategoryID: "",
			Amount:     decimal.NewFromFloat(50.00),
		}

		result, err := usecase.CreateTransaction(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionDTO{
			Name:       "Groceries",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(0.0),
		}

		result, err := usecase.CreateTransaction(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when category not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionDTO{
			Name:       "Groceries",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
		}

		mockIDGen.EXPECT().
			Generate().
			Return("txn-1")

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW: mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateTransaction(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully create transaction", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionDTO{
			Name:       "Groceries",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
		}

		category := &domain.Category{
			ID:     "cat-1",
			UserID: "user-1",
			Name:   "Food",
		}

		mockIDGen.EXPECT().
			Generate().
			Return("txn-1")

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(category, nil)

		mockTransactionRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txn *domain.Transaction) error {
				assert.Equal(t, "txn-1", txn.ID)
				assert.Equal(t, "user-1", txn.UserID)
				assert.Equal(t, "cat-1", txn.CategoryID)
				assert.Equal(t, decimal.NewFromFloat(50.00), txn.Amount)
				assert.Equal(t, "Groceries", txn.Name)
				assert.Nil(t, txn.RuleID)
				assert.Equal(t, fixedTime, txn.CreatedAt)
				assert.Equal(t, fixedTime, txn.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW:    mockCategoryRepo,
					TransactionRepositoryRW: mockTransactionRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateTransaction(ctx, "user-1", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "txn-1", result.ID)
		assert.Equal(t, "Groceries", result.Name)
	})
}

func TestTransactionUseCase_CreateTransactionRule(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepo := mocks.NewMockCategoryRepositoryRW(ctrl)
	mockRuleRepo := mocks.NewMockTransactionRuleRepositoryRW(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "Monthly rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			CronExpression: "0 0 1 * *",
		}

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "   ",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			CronExpression: "0 0 1 * *",
		}

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when cron_expression is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "Monthly rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			CronExpression: "   ",
		}

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when category_id is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "Monthly rent",
			CategoryID:     "",
			Amount:         decimal.NewFromFloat(1000.00),
			CronExpression: "0 0 1 * *",
		}

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "Monthly rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(0.0),
			CronExpression: "0 0 1 * *",
		}

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when category not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "Monthly rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			CronExpression: "0 0 1 * *",
		}

		mockIDGen.EXPECT().
			Generate().
			Return("rule-1")

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW: mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully create transaction rule", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateTransactionRuleDTO{
			Name:           "Monthly rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			CronExpression: "0 0 1 * *",
		}

		category := &domain.Category{
			ID:     "cat-1",
			UserID: "user-1",
			Name:   "Housing",
		}

		mockIDGen.EXPECT().
			Generate().
			Return("rule-1")

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(category, nil)

		mockRuleRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, rule *domain.TransactionRule) error {
				assert.Equal(t, "rule-1", rule.ID)
				assert.Equal(t, "user-1", rule.UserID)
				assert.Equal(t, "cat-1", rule.CategoryID)
				assert.Equal(t, decimal.NewFromFloat(1000.00), rule.Amount)
				assert.Equal(t, "Monthly rent", rule.Name)
				assert.Equal(t, "0 0 1 * *", rule.CronExpression)
				assert.Nil(t, rule.LastGeneratedAt)
				assert.Equal(t, fixedTime, rule.CreatedAt)
				assert.Equal(t, fixedTime, rule.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW:        mockCategoryRepo,
					TransactionRuleRepositoryRW: mockRuleRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateTransactionRule(ctx, "user-1", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "rule-1", result.ID)
		assert.Equal(t, "Monthly rent", result.Name)
		assert.Equal(t, "0 0 1 * *", result.CronExpression)
	})
}

func TestTransactionUseCase_GetTransaction(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockTransactionRepoR := mocks.NewMockTransactionRepositoryR(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot view resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		result, err := usecase.GetTransaction(ctx, "user-1", "txn-1")

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when transaction not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		mockTransactionRepoR.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "txn-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					TransactionRepositoryR: mockTransactionRepoR,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.GetTransaction(ctx, "user-1", "txn-1")

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully get transaction", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		transaction := &domain.Transaction{
			ID:         "txn-1",
			UserID:     "user-1",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
			Name:       "Groceries",
			RuleID:     nil,
			CreatedAt:  fixedTime,
			UpdatedAt:  fixedTime,
		}

		mockTransactionRepoR.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "txn-1").
			Return(transaction, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					TransactionRepositoryR: mockTransactionRepoR,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.GetTransaction(ctx, "user-1", "txn-1")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "txn-1", result.ID)
		assert.Equal(t, "Groceries", result.Name)
	})
}

func TestTransactionUseCase_GetTransactionRule(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockRuleRepoR := mocks.NewMockTransactionRuleRepositoryR(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot view resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		result, err := usecase.GetTransactionRule(ctx, "user-1", "rule-1")

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when transaction rule not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		mockRuleRepoR.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					TransactionRuleRepositoryR: mockRuleRepoR,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.GetTransactionRule(ctx, "user-1", "rule-1")

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully get transaction rule", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		rule := &domain.TransactionRule{
			ID:              "rule-1",
			UserID:          "user-1",
			CategoryID:      "cat-1",
			Amount:          decimal.NewFromFloat(1000.00),
			Name:            "Monthly rent",
			CronExpression:  "0 0 1 * *",
			LastGeneratedAt: nil,
			CreatedAt:       fixedTime,
			UpdatedAt:       fixedTime,
		}

		mockRuleRepoR.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(rule, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					TransactionRuleRepositoryR: mockRuleRepoR,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.GetTransactionRule(ctx, "user-1", "rule-1")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "rule-1", result.ID)
		assert.Equal(t, "Monthly rent", result.Name)
		assert.Equal(t, "0 0 1 * *", result.CronExpression)
	})
}

func TestTransactionUseCase_ListTransactions(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockTransactionRepoR := mocks.NewMockTransactionRepositoryR(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot view resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		result, total, err := usecase.ListTransactions(ctx, "user-1", nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, total)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should successfully list transactions", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		transactions := []*domain.Transaction{
			{
				ID:         "txn-1",
				UserID:     "user-1",
				CategoryID: "cat-1",
				Amount:     decimal.NewFromFloat(50.00),
				Name:       "Groceries",
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
			{
				ID:         "txn-2",
				UserID:     "user-1",
				CategoryID: "cat-2",
				Amount:     decimal.NewFromFloat(100.00),
				Name:       "Utilities",
				CreatedAt:  fixedTime,
				UpdatedAt:  fixedTime,
			},
		}

		mockTransactionRepoR.EXPECT().
			List(gomock.Any(), "user-1", nil, nil).
			Return(transactions, 2, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					TransactionRepositoryR: mockTransactionRepoR,
				}
				return fn(ctx, repo)
			})

		result, total, err := usecase.ListTransactions(ctx, "user-1", nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, total)
		assert.Len(t, result, 2)
		assert.Equal(t, "txn-1", result[0].ID)
		assert.Equal(t, "txn-2", result[1].ID)
	})
}

func TestTransactionUseCase_ListTransactionRules(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockRuleRepoR := mocks.NewMockTransactionRuleRepositoryR(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot view resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		result, total, err := usecase.ListTransactionRules(ctx, "user-1", nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, total)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should successfully list transaction rules", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		rules := []*domain.TransactionRule{
			{
				ID:             "rule-1",
				UserID:         "user-1",
				CategoryID:     "cat-1",
				Amount:         decimal.NewFromFloat(1000.00),
				Name:           "Monthly rent",
				CronExpression: "0 0 1 * *",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
			{
				ID:             "rule-2",
				UserID:         "user-1",
				CategoryID:     "cat-2",
				Amount:         decimal.NewFromFloat(50.00),
				Name:           "Weekly groceries",
				CronExpression: "0 0 * * 0",
				CreatedAt:      fixedTime,
				UpdatedAt:      fixedTime,
			},
		}

		mockRuleRepoR.EXPECT().
			ListEnabled(gomock.Any(), "user-1", nil, nil).
			Return(rules, 2, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					TransactionRuleRepositoryR: mockRuleRepoR,
				}
				return fn(ctx, repo)
			})

		result, total, err := usecase.ListTransactionRules(ctx, "user-1", nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, total)
		assert.Len(t, result, 2)
		assert.Equal(t, "rule-1", result[0].ID)
		assert.Equal(t, "rule-2", result[1].ID)
	})
}

func TestTransactionUseCase_UpdateTransaction(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepo := mocks.NewMockCategoryRepositoryRW(ctrl)
	mockTransactionRepo := mocks.NewMockTransactionRepositoryRW(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionDTO{
			Name:       "Updated Groceries",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(75.00),
		}

		result, err := usecase.UpdateTransaction(ctx, "user-1", "txn-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionDTO{
			Name:       "   ",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(75.00),
		}

		result, err := usecase.UpdateTransaction(ctx, "user-1", "txn-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when transaction not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionDTO{
			Name:       "Updated Groceries",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(75.00),
		}

		mockTransactionRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "txn-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRepositoryRW: mockTransactionRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateTransaction(ctx, "user-1", "txn-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return error when category not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionDTO{
			Name:       "Updated Groceries",
			CategoryID: "cat-2",
			Amount:     decimal.NewFromFloat(75.00),
		}

		transaction := &domain.Transaction{
			ID:         "txn-1",
			UserID:     "user-1",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
			Name:       "Groceries",
			CreatedAt:  fixedTime.Add(-24 * time.Hour),
			UpdatedAt:  fixedTime.Add(-24 * time.Hour),
		}

		mockTransactionRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "txn-1").
			Return(transaction, nil)

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-2").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRepositoryRW: mockTransactionRepo,
					CategoryRepositoryRW:    mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateTransaction(ctx, "user-1", "txn-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully update transaction", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionDTO{
			Name:       "Updated Groceries",
			CategoryID: "cat-2",
			Amount:     decimal.NewFromFloat(75.00),
		}

		transaction := &domain.Transaction{
			ID:         "txn-1",
			UserID:     "user-1",
			CategoryID: "cat-1",
			Amount:     decimal.NewFromFloat(50.00),
			Name:       "Groceries",
			CreatedAt:  fixedTime.Add(-24 * time.Hour),
			UpdatedAt:  fixedTime.Add(-24 * time.Hour),
		}

		category := &domain.Category{
			ID:     "cat-2",
			UserID: "user-1",
			Name:   "Food",
		}

		mockTransactionRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "txn-1").
			Return(transaction, nil)

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-2").
			Return(category, nil)

		mockTransactionRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, txn *domain.Transaction) error {
				assert.Equal(t, "txn-1", txn.ID)
				assert.Equal(t, "user-1", txn.UserID)
				assert.Equal(t, "cat-2", txn.CategoryID)
				assert.Equal(t, decimal.NewFromFloat(75.00), txn.Amount)
				assert.Equal(t, "Updated Groceries", txn.Name)
				assert.Equal(t, fixedTime, txn.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRepositoryRW: mockTransactionRepo,
					CategoryRepositoryRW:    mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateTransaction(ctx, "user-1", "txn-1", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "txn-1", result.ID)
		assert.Equal(t, "Updated Groceries", result.Name)
	})
}

func TestTransactionUseCase_UpdateTransactionRule(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepo := mocks.NewMockCategoryRepositoryRW(ctrl)
	mockRuleRepo := mocks.NewMockTransactionRuleRepositoryRW(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionRuleDTO{
			Name:           "Updated rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1200.00),
			CronExpression: "0 0 1 * *",
		}

		result, err := usecase.UpdateTransactionRule(ctx, "user-1", "rule-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionRuleDTO{
			Name:           "   ",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1200.00),
			CronExpression: "0 0 1 * *",
		}

		result, err := usecase.UpdateTransactionRule(ctx, "user-1", "rule-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when transaction rule not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionRuleDTO{
			Name:           "Updated rent",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1200.00),
			CronExpression: "0 0 1 * *",
		}

		mockRuleRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRuleRepositoryRW: mockRuleRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateTransactionRule(ctx, "user-1", "rule-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should return error when category not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionRuleDTO{
			Name:           "Updated rent",
			CategoryID:     "cat-2",
			Amount:         decimal.NewFromFloat(1200.00),
			CronExpression: "0 0 1 * *",
		}

		rule := &domain.TransactionRule{
			ID:             "rule-1",
			UserID:         "user-1",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			Name:           "Monthly rent",
			CronExpression: "0 0 1 * *",
			CreatedAt:      fixedTime.Add(-24 * time.Hour),
			UpdatedAt:      fixedTime.Add(-24 * time.Hour),
		}

		mockRuleRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(rule, nil)

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-2").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRuleRepositoryRW: mockRuleRepo,
					CategoryRepositoryRW:        mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateTransactionRule(ctx, "user-1", "rule-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully update transaction rule", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateTransactionRuleDTO{
			Name:           "Updated rent",
			CategoryID:     "cat-2",
			Amount:         decimal.NewFromFloat(1200.00),
			CronExpression: "0 0 2 * *",
		}

		rule := &domain.TransactionRule{
			ID:             "rule-1",
			UserID:         "user-1",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			Name:           "Monthly rent",
			CronExpression: "0 0 1 * *",
			CreatedAt:      fixedTime.Add(-24 * time.Hour),
			UpdatedAt:      fixedTime.Add(-24 * time.Hour),
		}

		category := &domain.Category{
			ID:     "cat-2",
			UserID: "user-1",
			Name:   "Housing",
		}

		mockRuleRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(rule, nil)

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-2").
			Return(category, nil)

		mockRuleRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, r *domain.TransactionRule) error {
				assert.Equal(t, "rule-1", r.ID)
				assert.Equal(t, "user-1", r.UserID)
				assert.Equal(t, "cat-2", r.CategoryID)
				assert.Equal(t, decimal.NewFromFloat(1200.00), r.Amount)
				assert.Equal(t, "Updated rent", r.Name)
				assert.Equal(t, "0 0 2 * *", r.CronExpression)
				assert.Equal(t, fixedTime, r.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRuleRepositoryRW: mockRuleRepo,
					CategoryRepositoryRW:        mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateTransactionRule(ctx, "user-1", "rule-1", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "rule-1", result.ID)
		assert.Equal(t, "Updated rent", result.Name)
		assert.Equal(t, "0 0 2 * *", result.CronExpression)
	})
}

func TestTransactionUseCase_DeleteTransactionRule(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockRuleRepo := mocks.NewMockTransactionRuleRepositoryRW(ctrl)

	usecase := usecases.NewTransactionUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		err := usecase.DeleteTransactionRule(ctx, "user-1", "rule-1")

		require.Error(t, err)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when transaction rule not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		mockRuleRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRuleRepositoryRW: mockRuleRepo,
				}
				return fn(ctx, repo)
			})

		err := usecase.DeleteTransactionRule(ctx, "user-1", "rule-1")

		require.Error(t, err)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully delete (disable) transaction rule", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		rule := &domain.TransactionRule{
			ID:             "rule-1",
			UserID:         "user-1",
			CategoryID:     "cat-1",
			Amount:         decimal.NewFromFloat(1000.00),
			Name:           "Monthly rent",
			CronExpression: "0 0 1 * *",
			DisabledAt:     nil,
			CreatedAt:      fixedTime.Add(-24 * time.Hour),
			UpdatedAt:      fixedTime.Add(-24 * time.Hour),
		}

		mockRuleRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "rule-1").
			Return(rule, nil)

		mockRuleRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, r *domain.TransactionRule) error {
				assert.Equal(t, "rule-1", r.ID)
				assert.Equal(t, "user-1", r.UserID)
				assert.NotNil(t, r.DisabledAt)
				assert.Equal(t, fixedTime, *r.DisabledAt)
				assert.Equal(t, fixedTime, r.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					TransactionRuleRepositoryRW: mockRuleRepo,
				}
				return fn(ctx, repo)
			})

		err := usecase.DeleteTransactionRule(ctx, "user-1", "rule-1")

		require.NoError(t, err)
	})
}
