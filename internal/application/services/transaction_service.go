package services

import (
	"cmp"
	"context"
	"fmt"
	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/domain"
	"time"

	"github.com/jonboulle/clockwork"
)

type Cron interface {
	GetNextGenerationTime(expression string, after time.Time) (time.Time, error)
}

type TransactionService struct {
	db          interfaces.Database
	clock       clockwork.Clock
	cron        Cron
	idGenerator interfaces.IDGenerator
}

func NewTransactionService(db interfaces.Database, clock clockwork.Clock, cron Cron, idGenerator interfaces.IDGenerator) *TransactionService {
	return &TransactionService{
		db:          db,
		clock:       clock,
		cron:        cron,
		idGenerator: idGenerator,
	}
}

func (s *TransactionService) GenerateTransactionFromRule(ctx context.Context, rule *domain.TransactionRule) error {
	now := s.clock.Now()

	after := cmp.Or(rule.LastGeneratedAt, &rule.CreatedAt)
	nextGenTime, err := s.cron.GetNextGenerationTime(rule.CronExpression, *after)
	if err != nil {
		return fmt.Errorf("failed to get next generation time: %w", err)
	}

	if nextGenTime.After(now) {
		return nil
	}

	return s.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		transactionID := s.idGenerator.Generate()
		err := repo.TransactionRepositoryRW.Create(ctx, rule.ToTransaction(transactionID, now))
		if err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		rule.LastGeneratedAt = &now
		return repo.TransactionRuleRepositoryRW.Update(ctx, rule)
	})
}

func (s *TransactionService) ClaimNextRuleToGenerate(ctx context.Context) (rule *domain.TransactionRule, err error) {
	err = s.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		rule, err = repo.TransactionRuleRepositoryRW.ClaimNextRuleToGenerate(ctx, s.clock.Now())
		return err
	})
	return rule, err
}
