package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"money_tracker/internal/domain"
)

type CategoryRepository struct {
	db pgx.Tx
}

func NewCategoryRepository(db pgx.Tx) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO categories (id, user_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		category.ID,
		category.UserID,
		category.Name,
		category.CreatedAt,
		category.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE categories
		SET user_id = $1, name = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.Exec(ctx, query,
		category.UserID,
		category.Name,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	return nil
}

func (r *CategoryRepository) GetByUserIDAndID(ctx context.Context, userID, categoryID string) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM categories
		WHERE user_id = $1 AND id = $2
	`

	var category domain.Category
	err := r.db.QueryRow(ctx, query, userID, categoryID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query category by user id and id: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) GetByUserIDAndName(ctx context.Context, userID, name string) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM categories
		WHERE user_id = $1 AND name = $2
	`

	var category domain.Category
	err := r.db.QueryRow(ctx, query, userID, name).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query category by user id and name: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) List(ctx context.Context, userID string, limit, offset *int) ([]*domain.Category, int, error) {
	countQuery := `SELECT COUNT(*) FROM categories WHERE user_id = $1`
	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count categories: %w", err)
	}

	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM categories
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
		return nil, 0, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate categories: %w", err)
	}

	return categories, total, nil
}
