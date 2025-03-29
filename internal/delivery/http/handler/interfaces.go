package handler

import (
	"context"

	usecaseDTO "money_tracker/internal/application/usecases/dto"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, dto usecaseDTO.CreateUserDTO) (*usecaseDTO.UserDTO, error)
}

type AuthUseCase interface {
	Login(ctx context.Context, dto usecaseDTO.LoginDTO) (*usecaseDTO.AuthResponseDTO, error)
}

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, userID string, dto usecaseDTO.CreateTransactionDTO) (*usecaseDTO.TransactionDTO, error)
	GetTransaction(ctx context.Context, userID string, transactionID string) (*usecaseDTO.TransactionDTO, error)
	ListTransactions(ctx context.Context, userID string, limit, offset *int) ([]*usecaseDTO.TransactionDTO, int, error)
	UpdateTransaction(ctx context.Context, userID string, transactionID string, dto usecaseDTO.UpdateTransactionDTO) (*usecaseDTO.TransactionDTO, error)
}

type CategoryUseCase interface {
	CreateCategory(ctx context.Context, userID string, dto usecaseDTO.CreateCategoryDTO) (*usecaseDTO.CategoryDTO, error)
	GetCategory(ctx context.Context, userID string, categoryID string) (*usecaseDTO.CategoryDTO, error)
	ListCategories(ctx context.Context, userID string, limit, offset *int) ([]*usecaseDTO.CategoryDTO, int, error)
	UpdateCategory(ctx context.Context, userID string, categoryID string, dto usecaseDTO.UpdateCategoryDTO) (*usecaseDTO.CategoryDTO, error)
}

type TransactionRuleUseCase interface {
	CreateTransactionRule(ctx context.Context, userID string, dto usecaseDTO.CreateTransactionRuleDTO) (*usecaseDTO.TransactionRuleDTO, error)
	GetTransactionRule(ctx context.Context, userID string, ruleID string) (*usecaseDTO.TransactionRuleDTO, error)
	ListTransactionRules(ctx context.Context, userID string, limit, offset *int) ([]*usecaseDTO.TransactionRuleDTO, int, error)
	UpdateTransactionRule(ctx context.Context, userID string, ruleID string, dto usecaseDTO.UpdateTransactionRuleDTO) (*usecaseDTO.TransactionRuleDTO, error)
	DeleteTransactionRule(ctx context.Context, userID string, ruleID string) error
}
