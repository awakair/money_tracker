package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"money_tracker/internal/application/interfaces"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

func (db *Database) withTransaction(ctx context.Context, txOptions pgx.TxOptions, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := db.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (db *Database) DoReadWrite(ctx context.Context, fn func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error) error {
	return db.withTransaction(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	}, func(ctx context.Context, tx pgx.Tx) error {
		repo := interfaces.RepositoriesReadWrite{
			UserRepositoryRW:            NewUserRepository(tx),
			CategoryRepositoryRW:        NewCategoryRepository(tx),
			TransactionRepositoryRW:     NewTransactionRepository(tx),
			TransactionRuleRepositoryRW: NewTransactionRuleRepository(tx),
		}
		return fn(ctx, repo)
	})
}

func (db *Database) DoRead(ctx context.Context, fn func(ctx context.Context, repo interfaces.RepositoriesRead) error) error {
	return db.withTransaction(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	}, func(ctx context.Context, tx pgx.Tx) error {
		repo := interfaces.RepositoriesRead{
			UserRepositoryR:            NewUserRepository(tx),
			CategoryRepositoryR:        NewCategoryRepository(tx),
			TransactionRepositoryR:     NewTransactionRepository(tx),
			TransactionRuleRepositoryR: NewTransactionRuleRepository(tx),
		}
		return fn(ctx, repo)
	})
}
