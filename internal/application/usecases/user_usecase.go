package usecases

import (
	"context"
	"fmt"

	"github.com/jonboulle/clockwork"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
)

type UserUseCase struct {
	db             interfaces.Database
	passwordHasher interfaces.PasswordHasher
	idGenerator    interfaces.IDGenerator
	clock          clockwork.Clock
}

func NewUserUseCase(
	db interfaces.Database,
	passwordHasher interfaces.PasswordHasher,
	idGenerator interfaces.IDGenerator,
	clock clockwork.Clock,
) *UserUseCase {
	return &UserUseCase{
		db:             db,
		passwordHasher: passwordHasher,
		idGenerator:    idGenerator,
		clock:          clock,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, input dto.CreateUserDTO) (*dto.UserDTO, error) {
	if input.Email == "" {
		return nil, &domain.InvalidArgumentError{Argument: "email", Reason: "email is required"}
	}

	if len(input.Password) < 8 {
		return nil, &domain.InvalidArgumentError{Argument: "password", Reason: "must be at least 8 characters"}
	}

	passwordHash, err := uc.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := uc.clock.Now()
	user := &domain.User{
		ID:           uc.idGenerator.Generate(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         domain.UserRoleUser,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		existingUser, err := repo.UserRepositoryRW.GetByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("check for existing user: %w", err)
		}

		if existingUser != nil {
			return &domain.EntityAlreadyExistsError{EntityType: "User", Field: "email", Value: input.Email}
		}

		if err := repo.UserRepositoryRW.Create(ctx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromUser(user), nil
}
