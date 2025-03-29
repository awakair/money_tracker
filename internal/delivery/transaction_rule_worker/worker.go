package transaction_rule_worker

import (
	"context"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"

	"money_tracker/internal/domain"
	"money_tracker/internal/pkg/logger"
)

type TransactionService interface {
	ClaimNextRuleToGenerate(ctx context.Context) (*domain.TransactionRule, error)
	GenerateTransactionFromRule(ctx context.Context, rule *domain.TransactionRule) error
}

type Worker struct {
	service      TransactionService
	logger       logger.Logger
	clock        clockwork.Clock
	pollInterval time.Duration
}

func NewWorker(service TransactionService, logger logger.Logger, clock clockwork.Clock, pollInterval time.Duration) *Worker {
	return &Worker{
		service:      service,
		logger:       logger,
		clock:        clock,
		pollInterval: pollInterval,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.logger.InfoContext(ctx, "Transaction rule worker started", "poll_interval", w.pollInterval.String())

	ticker := w.clock.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.InfoContext(ctx, "Transaction rule worker stopped")
			return ctx.Err()
		case <-ticker.Chan():
			if err := w.processNextRule(ctx); err != nil {
				w.logger.ErrorContext(ctx, "Failed to process rule", "error", err.Error())
			}
		}
	}
}

func (w *Worker) processNextRule(ctx context.Context) error {
	rule, err := w.service.ClaimNextRuleToGenerate(ctx)
	if err != nil {
		return fmt.Errorf("claim next rule: %w", err)
	}

	if rule == nil {
		return nil
	}

	if err := w.service.GenerateTransactionFromRule(ctx, rule); err != nil {
		return fmt.Errorf("generate transaction from rule: %w", err)
	}

	w.logger.InfoContext(ctx, "Successfully ran worker iteration", "rule_id", rule.ID)
	return nil
}
