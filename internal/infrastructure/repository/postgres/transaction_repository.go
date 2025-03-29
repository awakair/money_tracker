package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"money_tracker/internal/domain"
)

type TransactionRepository struct {
	db pgx.Tx
}

func NewTransactionRepository(db pgx.Tx) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, category_id, amount, name, rule_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		transaction.ID,
		transaction.UserID,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Name,
		transaction.RuleID,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET user_id = $1, category_id = $2, amount = $3, name = $4, rule_id = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(ctx, query,
		transaction.UserID,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Name,
		transaction.RuleID,
		transaction.UpdatedAt,
		transaction.ID,
	)
	if err != nil {
		return fmt.Errorf("update transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetByUserIDAndID(ctx context.Context, userID, transactionID string) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, category_id, amount, name, rule_id, created_at, updated_at
		FROM transactions
		WHERE user_id = $1 AND id = $2
	`

	var transaction domain.Transaction
	err := r.db.QueryRow(ctx, query, userID, transactionID).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.CategoryID,
		&transaction.Amount,
		&transaction.Name,
		&transaction.RuleID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query transaction by user id and id: %w", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) List(ctx context.Context, userID string, limit, offset *int) ([]*domain.Transaction, int, error) {
	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count transactions: %w", err)
	}

	query := `
		SELECT id, user_id, category_id, amount, name, rule_id, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
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
		return nil, 0, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.CategoryID,
			&transaction.Amount,
			&transaction.Name,
			&transaction.RuleID,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan transaction: %w", err)
		}
		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate transactions: %w", err)
	}

	return transactions, total, nil
}
