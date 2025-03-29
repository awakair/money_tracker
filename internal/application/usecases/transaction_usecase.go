package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/jonboulle/clockwork"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
	"money_tracker/internal/pkg/credscontext"
)

type TransactionUseCase struct {
	db          interfaces.Database
	idGenerator interfaces.IDGenerator
	clock       clockwork.Clock
}

func NewTransactionUseCase(
	db interfaces.Database,
	idGenerator interfaces.IDGenerator,
	clock clockwork.Clock,
) *TransactionUseCase {
	return &TransactionUseCase{
		db:          db,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, userID string, input dto.CreateTransactionDTO) (*dto.TransactionDTO, error) {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, &domain.InvalidArgumentError{Argument: "name", Reason: "name is required"}
	}

	if input.CategoryID == "" {
		return nil, &domain.InvalidArgumentError{Argument: "category_id", Reason: "category_id is required"}
	}

	if !input.Amount.IsPositive() {
		return nil, &domain.InvalidArgumentError{Argument: "amount", Reason: "amount must be positive"}
	}

	now := uc.clock.Now()
	transaction := &domain.Transaction{
		ID:         uc.idGenerator.Generate(),
		UserID:     userID,
		CategoryID: input.CategoryID,
		Amount:     input.Amount,
		Name:       input.Name,
		RuleID:     nil,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		category, err := repo.CategoryRepositoryRW.GetByUserIDAndID(ctx, userID, input.CategoryID)
		if err != nil {
			return fmt.Errorf("get category: %w", err)
		}

		if category == nil {
			return &domain.EntityNotFoundError{EntityType: "Category", ID: input.CategoryID}
		}

		if err := repo.TransactionRepositoryRW.Create(ctx, transaction); err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromTransaction(transaction), nil
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, userID, transactionID string) (*dto.TransactionDTO, error) {
	if !domain.CanViewResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	var transaction *domain.Transaction
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		transaction, err = repo.TransactionRepositoryR.GetByUserIDAndID(ctx, userID, transactionID)
		if err != nil {
			return fmt.Errorf("get transaction: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, &domain.EntityNotFoundError{EntityType: "Transaction", ID: transactionID}
	}

	return dto.FromTransaction(transaction), nil
}

func (uc *TransactionUseCase) ListTransactions(ctx context.Context, userID string, limit, offset *int) ([]*dto.TransactionDTO, int, error) {
	if !domain.CanViewResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, 0, &domain.AuthorizationError{}
	}

	var transactions []*domain.Transaction
	var total int
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		transactions, total, err = repo.TransactionRepositoryR.List(ctx, userID, limit, offset)
		if err != nil {
			return fmt.Errorf("list transactions: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.TransactionDTO, len(transactions))
	for i, txn := range transactions {
		result[i] = dto.FromTransaction(txn)
	}

	return result, total, nil
}

func (uc *TransactionUseCase) UpdateTransaction(ctx context.Context, userID, transactionID string, input dto.UpdateTransactionDTO) (*dto.TransactionDTO, error) {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, &domain.InvalidArgumentError{Argument: "name", Reason: "name is required"}
	}

	if input.CategoryID == "" {
		return nil, &domain.InvalidArgumentError{Argument: "category_id", Reason: "category_id is required"}
	}

	if !input.Amount.IsPositive() {
		return nil, &domain.InvalidArgumentError{Argument: "amount", Reason: "amount must be positive"}
	}

	var transaction *domain.Transaction
	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		var err error
		transaction, err = repo.TransactionRepositoryRW.GetByUserIDAndID(ctx, userID, transactionID)
		if err != nil {
			return fmt.Errorf("get transaction: %w", err)
		}

		if transaction == nil {
			return &domain.EntityNotFoundError{EntityType: "Transaction", ID: transactionID}
		}

		category, err := repo.CategoryRepositoryRW.GetByUserIDAndID(ctx, userID, input.CategoryID)
		if err != nil {
			return fmt.Errorf("get category: %w", err)
		}

		if category == nil {
			return &domain.EntityNotFoundError{EntityType: "Category", ID: input.CategoryID}
		}

		transaction.Name = input.Name
		transaction.CategoryID = input.CategoryID
		transaction.Amount = input.Amount
		transaction.UpdatedAt = uc.clock.Now()

		if err := repo.TransactionRepositoryRW.Update(ctx, transaction); err != nil {
			return fmt.Errorf("update transaction: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromTransaction(transaction), nil
}

func (uc *TransactionUseCase) CreateTransactionRule(ctx context.Context, userID string, input dto.CreateTransactionRuleDTO) (*dto.TransactionRuleDTO, error) {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	input.Name = strings.TrimSpace(input.Name)
	input.CronExpression = strings.TrimSpace(input.CronExpression)

	if input.Name == "" {
		return nil, &domain.InvalidArgumentError{Argument: "name", Reason: "name is required"}
	}

	if input.CategoryID == "" {
		return nil, &domain.InvalidArgumentError{Argument: "category_id", Reason: "category_id is required"}
	}

	if !input.Amount.IsPositive() {
		return nil, &domain.InvalidArgumentError{Argument: "amount", Reason: "amount must be positive"}
	}

	if input.CronExpression == "" {
		return nil, &domain.InvalidArgumentError{Argument: "cron_expression", Reason: "cron_expression is required"}
	}

	now := uc.clock.Now()
	rule := &domain.TransactionRule{
		ID:              uc.idGenerator.Generate(),
		UserID:          userID,
		CategoryID:      input.CategoryID,
		Amount:          input.Amount,
		Name:            input.Name,
		CronExpression:  input.CronExpression,
		LastGeneratedAt: nil,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		category, err := repo.CategoryRepositoryRW.GetByUserIDAndID(ctx, userID, input.CategoryID)
		if err != nil {
			return fmt.Errorf("get category: %w", err)
		}

		if category == nil {
			return &domain.EntityNotFoundError{EntityType: "Category", ID: input.CategoryID}
		}

		if err := repo.TransactionRuleRepositoryRW.Create(ctx, rule); err != nil {
			return fmt.Errorf("create transaction rule: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromTransactionRule(rule), nil
}

func (uc *TransactionUseCase) GetTransactionRule(ctx context.Context, userID, ruleID string) (*dto.TransactionRuleDTO, error) {
	if !domain.CanViewResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	var rule *domain.TransactionRule
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		rule, err = repo.TransactionRuleRepositoryR.GetByUserIDAndID(ctx, userID, ruleID)
		if err != nil {
			return fmt.Errorf("get transaction rule: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, &domain.EntityNotFoundError{EntityType: "TransactionRule", ID: ruleID}
	}

	return dto.FromTransactionRule(rule), nil
}

func (uc *TransactionUseCase) ListTransactionRules(ctx context.Context, userID string, limit, offset *int) ([]*dto.TransactionRuleDTO, int, error) {
	if !domain.CanViewResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, 0, &domain.AuthorizationError{}
	}

	var rules []*domain.TransactionRule
	var total int
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		rules, total, err = repo.TransactionRuleRepositoryR.ListEnabled(ctx, userID, limit, offset)
		if err != nil {
			return fmt.Errorf("list transaction rules: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.TransactionRuleDTO, len(rules))
	for i, rule := range rules {
		result[i] = dto.FromTransactionRule(rule)
	}

	return result, total, nil
}

func (uc *TransactionUseCase) UpdateTransactionRule(ctx context.Context, userID, ruleID string, input dto.UpdateTransactionRuleDTO) (*dto.TransactionRuleDTO, error) {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	input.Name = strings.TrimSpace(input.Name)
	input.CronExpression = strings.TrimSpace(input.CronExpression)

	if input.Name == "" {
		return nil, &domain.InvalidArgumentError{Argument: "name", Reason: "name is required"}
	}

	if input.CategoryID == "" {
		return nil, &domain.InvalidArgumentError{Argument: "category_id", Reason: "category_id is required"}
	}

	if !input.Amount.IsPositive() {
		return nil, &domain.InvalidArgumentError{Argument: "amount", Reason: "amount must be positive"}
	}

	if input.CronExpression == "" {
		return nil, &domain.InvalidArgumentError{Argument: "cron_expression", Reason: "cron_expression is required"}
	}

	var rule *domain.TransactionRule
	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		var err error
		rule, err = repo.TransactionRuleRepositoryRW.GetByUserIDAndID(ctx, userID, ruleID)
		if err != nil {
			return fmt.Errorf("get transaction rule: %w", err)
		}

		if rule == nil {
			return &domain.EntityNotFoundError{EntityType: "TransactionRule", ID: ruleID}
		}

		category, err := repo.CategoryRepositoryRW.GetByUserIDAndID(ctx, userID, input.CategoryID)
		if err != nil {
			return fmt.Errorf("get category: %w", err)
		}

		if category == nil {
			return &domain.EntityNotFoundError{EntityType: "Category", ID: input.CategoryID}
		}

		rule.Name = input.Name
		rule.CategoryID = input.CategoryID
		rule.Amount = input.Amount
		rule.CronExpression = input.CronExpression
		rule.UpdatedAt = uc.clock.Now()

		if err := repo.TransactionRuleRepositoryRW.Update(ctx, rule); err != nil {
			return fmt.Errorf("update transaction rule: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromTransactionRule(rule), nil
}

func (uc *TransactionUseCase) DeleteTransactionRule(ctx context.Context, userID, ruleID string) error {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return &domain.AuthorizationError{}
	}

	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		rule, err := repo.TransactionRuleRepositoryRW.GetByUserIDAndID(ctx, userID, ruleID)
		if err != nil {
			return fmt.Errorf("get transaction rule: %w", err)
		}

		if rule == nil {
			return &domain.EntityNotFoundError{EntityType: "TransactionRule", ID: ruleID}
		}

		now := uc.clock.Now()
		rule.DisabledAt = &now
		rule.UpdatedAt = now

		if err := repo.TransactionRuleRepositoryRW.Update(ctx, rule); err != nil {
			return fmt.Errorf("disable transaction rule: %w", err)
		}

		return nil
	})

	return err
}
