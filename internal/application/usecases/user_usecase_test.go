package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/interfaces/mocks"
	"money_tracker/internal/application/usecases"
	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockHasher := mocks.NewMockPasswordHasher(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockUserRepo := mocks.NewMockUserRepositoryRW(ctrl)

	usecase := usecases.NewUserUseCase(mockDB, mockHasher, mockIDGen, clock)

	t.Run("should return error when email is empty", func(t *testing.T) {
		input := dto.CreateUserDTO{
			Email:    "",
			Password: "password123",
		}

		result, err := usecase.CreateUser(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when password is too short", func(t *testing.T) {
		input := dto.CreateUserDTO{
			Email:    "test@example.com",
			Password: "short",
		}

		result, err := usecase.CreateUser(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when user with email already exists", func(t *testing.T) {
		input := dto.CreateUserDTO{
			Email:    "test@example.com",
			Password: "password123",
		}

		existingUser := &domain.User{
			ID:    "user-1",
			Email: "test@example.com",
		}

		mockHasher.EXPECT().
			Hash("password123").
			Return("hashed_password", nil)

		mockIDGen.EXPECT().
			Generate().
			Return("user-new")

		mockUserRepo.EXPECT().
			GetByEmail(gomock.Any(), "test@example.com").
			Return(existingUser, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					UserRepositoryRW: mockUserRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateUser(ctx, input)

		require.Error(t, err)
		assert.Nil(t, result)

		var entityExistsErr *domain.EntityAlreadyExistsError
		require.ErrorAs(t, err, &entityExistsErr)
	})

	t.Run("should successfully create user", func(t *testing.T) {
		input := dto.CreateUserDTO{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockHasher.EXPECT().
			Hash("password123").
			Return("hashed_password", nil)

		mockIDGen.EXPECT().
			Generate().
			Return("user-1")

		mockUserRepo.EXPECT().
			GetByEmail(gomock.Any(), "test@example.com").
			Return(nil, nil)

		mockUserRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *domain.User) error {
				assert.Equal(t, "user-1", user.ID)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "hashed_password", user.PasswordHash)
				assert.Equal(t, domain.UserRoleUser, user.Role)
				assert.Equal(t, fixedTime, user.CreatedAt)
				assert.Equal(t, fixedTime, user.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					UserRepositoryRW: mockUserRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateUser(ctx, input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "user-1", result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, string(domain.UserRoleUser), result.Role)
	})
}
