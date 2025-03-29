package interfaces

import (
	"context"
	"money_tracker/internal/domain"
	"time"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

type IDGenerator interface {
	Generate() string
}

type Database interface {
	DoReadWrite(ctx context.Context, fn func(ctx context.Context, repo RepositoriesReadWrite) error) error
	DoRead(ctx context.Context, fn func(ctx context.Context, repo RepositoriesRead) error) error
}

type RepositoriesReadWrite struct {
	UserRepositoryRW
	CategoryRepositoryRW
	TransactionRepositoryRW
	TransactionRuleRepositoryRW
}

type RepositoriesRead struct {
	UserRepositoryR
	CategoryRepositoryR
	TransactionRepositoryR
	TransactionRuleRepositoryR
}

type UserRepositoryRW interface {
	UserRepositoryR
	Create(ctx context.Context, user *domain.User) error
}

type UserRepositoryR interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type CategoryRepositoryRW interface {
	CategoryRepositoryR
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
}

type CategoryRepositoryR interface {
	GetByUserIDAndID(ctx context.Context, userID, categoryID string) (*domain.Category, error)
	GetByUserIDAndName(ctx context.Context, userID, name string) (*domain.Category, error)
	List(ctx context.Context, userID string, limit, offset *int) ([]*domain.Category, int, error)
}

type TransactionRepositoryRW interface {
	TransactionRepositoryR
	Create(ctx context.Context, transaction *domain.Transaction) error
	Update(ctx context.Context, transaction *domain.Transaction) error
}

type TransactionRepositoryR interface {
	GetByUserIDAndID(ctx context.Context, userID, transactionID string) (*domain.Transaction, error)
	List(ctx context.Context, userID string, limit, offset *int) ([]*domain.Transaction, int, error)
}

type TransactionRuleRepositoryRW interface {
	TransactionRuleRepositoryR
	Create(ctx context.Context, rule *domain.TransactionRule) error
	Update(ctx context.Context, rule *domain.TransactionRule) error
	ClaimNextRuleToGenerate(ctx context.Context, now time.Time) (*domain.TransactionRule, error)
}

type TransactionRuleRepositoryR interface {
	GetByUserIDAndID(ctx context.Context, userID, ruleID string) (*domain.TransactionRule, error)
	ListEnabled(ctx context.Context, userID string, limit, offset *int) ([]*domain.TransactionRule, int, error)
}
