package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"money_tracker/internal/domain"
)

type TransactionRuleRepository struct {
	db pgx.Tx
}

func NewTransactionRuleRepository(db pgx.Tx) *TransactionRuleRepository {
	return &TransactionRuleRepository{db: db}
}

func (r *TransactionRuleRepository) Create(ctx context.Context, rule *domain.TransactionRule) error {
	query := `
		INSERT INTO transaction_rules (id, user_id, category_id, amount, name, cron_expression, last_generated_at, disabled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query,
		rule.ID,
		rule.UserID,
		rule.CategoryID,
		rule.Amount,
		rule.Name,
		rule.CronExpression,
		rule.LastGeneratedAt,
		rule.DisabledAt,
		rule.CreatedAt,
		rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert transaction rule: %w", err)
	}

	return nil
}

func (r *TransactionRuleRepository) Update(ctx context.Context, rule *domain.TransactionRule) error {
	query := `
		UPDATE transaction_rules
		SET user_id = $1, category_id = $2, amount = $3, name = $4, cron_expression = $5, last_generated_at = $6, disabled_at = $7, updated_at = $8
		WHERE id = $9
	`

	_, err := r.db.Exec(ctx, query,
		rule.UserID,
		rule.CategoryID,
		rule.Amount,
		rule.Name,
		rule.CronExpression,
		rule.LastGeneratedAt,
		rule.DisabledAt,
		rule.UpdatedAt,
		rule.ID,
	)
	if err != nil {
		return fmt.Errorf("update transaction rule: %w", err)
	}

	return nil
}

func (r *TransactionRuleRepository) GetByUserIDAndID(ctx context.Context, userID, ruleID string) (*domain.TransactionRule, error) {
	query := `
		SELECT id, user_id, category_id, amount, name, cron_expression, last_generated_at, disabled_at, created_at, updated_at
		FROM transaction_rules
		WHERE user_id = $1 AND id = $2 AND disabled_at IS NULL
	`

	var rule domain.TransactionRule
	err := r.db.QueryRow(ctx, query, userID, ruleID).Scan(
		&rule.ID,
		&rule.UserID,
		&rule.CategoryID,
		&rule.Amount,
		&rule.Name,
		&rule.CronExpression,
		&rule.LastGeneratedAt,
		&rule.DisabledAt,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query transaction rule by user id and id: %w", err)
	}

	return &rule, nil
}

func (r *TransactionRuleRepository) ListEnabled(ctx context.Context, userID string, limit, offset *int) ([]*domain.TransactionRule, int, error) {
	countQuery := `SELECT COUNT(*) FROM transaction_rules WHERE user_id = $1 AND disabled_at IS NULL`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count transaction rules: %w", err)
	}

	query := `
		SELECT id, user_id, category_id, amount, name, cron_expression, last_generated_at, disabled_at, created_at, updated_at
		FROM transaction_rules
		WHERE user_id = $1 AND disabled_at IS NULL
		ORDER BY created_at DESC
	`

	if limit != nil {
		query += fmt.Sprintf(` LIMIT %d`, *limit)
	}

	if offset != nil {
		query += fmt.Sprintf(` OFFSET %d`, *offset)
	}

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("query transaction rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.TransactionRule
	for rows.Next() {
		var rule domain.TransactionRule
		err := rows.Scan(
			&rule.ID,
			&rule.UserID,
			&rule.CategoryID,
			&rule.Amount,
			&rule.Name,
			&rule.CronExpression,
			&rule.LastGeneratedAt,
			&rule.DisabledAt,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction rule: %w", err)
		}
		rules = append(rules, &rule)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate transaction rules: %w", err)
	}

	return rules, total, nil
}

func (r *TransactionRuleRepository) ClaimNextRuleToGenerate(ctx context.Context, now time.Time) (*domain.TransactionRule, error) {
	query := `
		UPDATE transaction_rules
		SET updated_at = $1
		WHERE id = (
			SELECT id
			FROM transaction_rules
			WHERE disabled_at IS NULL
			ORDER BY updated_at ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, user_id, category_id, amount, name, cron_expression, last_generated_at, disabled_at, created_at, updated_at
	`

	var rule domain.TransactionRule
	err := r.db.QueryRow(ctx, query, now).Scan(
		&rule.ID,
		&rule.UserID,
		&rule.CategoryID,
		&rule.Amount,
		&rule.Name,
		&rule.CronExpression,
		&rule.LastGeneratedAt,
		&rule.DisabledAt,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("claim next rule to generate: %w", err)
	}

	return &rule, nil
}
