package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/interfaces/mocks"
	"money_tracker/internal/application/usecases"
	"money_tracker/internal/application/usecases/dto"
	usecasemocks "money_tracker/internal/application/usecases/mocks"
	"money_tracker/internal/domain"
)

func TestAuthUseCase_Login(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockHasher := mocks.NewMockPasswordHasher(ctrl)
	mockTokenGen := usecasemocks.NewMockTokenGenerator(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryR(ctrl)

	usecase := usecases.NewAuthUseCase(mockDB, mockHasher, mockTokenGen)

	t.Run("should return error when email is empty", func(t *testing.T) {
		input := dto.LoginDTO{
			Email:    "",
			Password: "password123",
		}

		result, err := usecase.Login(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		input := dto.LoginDTO{
			Email:    "test@example.com",
			Password: "",
		}

		result, err := usecase.Login(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return authorization error when user not found", func(t *testing.T) {
		input := dto.LoginDTO{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserRepo.EXPECT().
			GetByEmail(gomock.Any(), "test@example.com").
			Return(nil, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					UserRepositoryR: mockUserRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.Login(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return authorization error when password is incorrect", func(t *testing.T) {
		input := dto.LoginDTO{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		user := &domain.User{
			ID:           "user-1",
			Email:        "test@example.com",
			PasswordHash: "hashed_password",
			Role:         domain.UserRoleUser,
		}

		mockUserRepo.EXPECT().
			GetByEmail(gomock.Any(), "test@example.com").
			Return(user, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					UserRepositoryR: mockUserRepo,
				}
				return fn(ctx, repo)
			})

		mockHasher.EXPECT().
			Verify("wrongpassword", "hashed_password").
			Return(errors.New("password mismatch"))

		result, err := usecase.Login(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should successfully login with correct credentials", func(t *testing.T) {
		input := dto.LoginDTO{
			Email:    "test@example.com",
			Password: "password123",
		}

		user := &domain.User{
			ID:           "user-1",
			Email:        "test@example.com",
			PasswordHash: "hashed_password",
			Role:         domain.UserRoleUser,
		}

		mockUserRepo.EXPECT().
			GetByEmail(gomock.Any(), "test@example.com").
			Return(user, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					UserRepositoryR: mockUserRepo,
				}
				return fn(ctx, repo)
			})

		mockHasher.EXPECT().
			Verify("password123", "hashed_password").
			Return(nil)

		mockTokenGen.EXPECT().
			Generate("user-1", domain.UserRoleUser).
			Return("jwt_token_123", nil)

		result, err := usecase.Login(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "jwt_token_123", result.Token)
		assert.Equal(t, "user-1", result.User.ID)
		assert.Equal(t, "test@example.com", result.User.Email)
		assert.Equal(t, string(domain.UserRoleUser), result.User.Role)
	})
}
