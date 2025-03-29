package usecases

import (
	"context"
	"fmt"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
)

type TokenGenerator interface {
	Generate(userID string, role domain.UserRole) (string, error)
}

type AuthUseCase struct {
	db             interfaces.Database
	passwordHasher interfaces.PasswordHasher
	tokenGenerator TokenGenerator
}

func NewAuthUseCase(
	db interfaces.Database,
	passwordHasher interfaces.PasswordHasher,
	tokenGenerator TokenGenerator,
) *AuthUseCase {
	return &AuthUseCase{
		db:             db,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, input dto.LoginDTO) (*dto.AuthResponseDTO, error) {
	if input.Email == "" {
		return nil, &domain.InvalidArgumentError{Argument: "email", Reason: "email is required"}
	}

	if input.Password == "" {
		return nil, &domain.InvalidArgumentError{Argument: "password", Reason: "password is required"}
	}

	var user *domain.User
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		user, err = repo.UserRepositoryR.GetByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("get user by email: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &domain.AuthorizationError{}
	}

	if err := uc.passwordHasher.Verify(input.Password, user.PasswordHash); err != nil {
		return nil, &domain.AuthorizationError{}
	}

	token, err := uc.tokenGenerator.Generate(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &dto.AuthResponseDTO{
		Token: token,
		User:  dto.FromUser(user),
	}, nil
}
