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
	"money_tracker/internal/pkg/credscontext"
)

func TestCategoryUseCase_CreateCategory(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepo := mocks.NewMockCategoryRepositoryRW(ctrl)

	usecase := usecases.NewCategoryUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateCategoryDTO{
			Name: "Food",
		}

		result, err := usecase.CreateCategory(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateCategoryDTO{
			Name: "   ",
		}

		result, err := usecase.CreateCategory(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when category with same name already exists", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateCategoryDTO{
			Name: "  Food",
		}

		existingCategory := &domain.Category{
			ID:     "cat-1",
			UserID: "user-1",
			Name:   "Food",
		}

		mockIDGen.EXPECT().
			Generate().
			Return("cat-new")

		mockCategoryRepo.EXPECT().
			GetByUserIDAndName(gomock.Any(), "user-1", "Food").
			Return(existingCategory, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW: mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateCategory(ctx, "user-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var entityExistsErr *domain.EntityAlreadyExistsError
		require.ErrorAs(t, err, &entityExistsErr)
	})

	t.Run("should successfully create category", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.CreateCategoryDTO{
			Name: "Food",
		}

		mockIDGen.EXPECT().
			Generate().
			Return("cat-1")

		mockCategoryRepo.EXPECT().
			GetByUserIDAndName(gomock.Any(), "user-1", "Food").
			Return(nil, nil)

		mockCategoryRepo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cat *domain.Category) error {
				assert.Equal(t, "cat-1", cat.ID)
				assert.Equal(t, "user-1", cat.UserID)
				assert.Equal(t, "Food", cat.Name)
				assert.Equal(t, fixedTime, cat.CreatedAt)
				assert.Equal(t, fixedTime, cat.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW: mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.CreateCategory(ctx, "user-1", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "cat-1", result.ID)
		assert.Equal(t, "Food", result.Name)
	})
}

func TestCategoryUseCase_UpdateCategory(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepo := mocks.NewMockCategoryRepositoryRW(ctrl)

	usecase := usecases.NewCategoryUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot edit resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateCategoryDTO{
			Name: "Updated Food",
		}

		result, err := usecase.UpdateCategory(ctx, "user-1", "cat-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateCategoryDTO{
			Name: "   ",
		}

		result, err := usecase.UpdateCategory(ctx, "user-1", "cat-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var invalidArgErr *domain.InvalidArgumentError
		require.ErrorAs(t, err, &invalidArgErr)
	})

	t.Run("should return error when category not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateCategoryDTO{
			Name: "Updated Food",
		}

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW: mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateCategory(ctx, "user-1", "cat-1", input)

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully update category", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		input := dto.UpdateCategoryDTO{
			Name: "Updated Food",
		}

		category := &domain.Category{
			ID:        "cat-1",
			UserID:    "user-1",
			Name:      "Food",
			CreatedAt: fixedTime.Add(-24 * time.Hour),
			UpdatedAt: fixedTime.Add(-24 * time.Hour),
		}

		mockCategoryRepo.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(category, nil)

		mockCategoryRepo.EXPECT().
			GetByUserIDAndName(gomock.Any(), "user-1", "Updated Food").
			Return(nil, nil)

		mockCategoryRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cat *domain.Category) error {
				assert.Equal(t, "cat-1", cat.ID)
				assert.Equal(t, "user-1", cat.UserID)
				assert.Equal(t, "Updated Food", cat.Name)
				assert.Equal(t, fixedTime, cat.UpdatedAt)
				return nil
			})

		mockDB.EXPECT().
			DoReadWrite(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesReadWrite) error) error {
				repo := interfaces.RepositoriesReadWrite{
					CategoryRepositoryRW: mockCategoryRepo,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.UpdateCategory(ctx, "user-1", "cat-1", input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "cat-1", result.ID)
		assert.Equal(t, "Updated Food", result.Name)
	})
}

func TestCategoryUseCase_GetCategory(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepoR := mocks.NewMockCategoryRepositoryR(ctrl)

	usecase := usecases.NewCategoryUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot view resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		result, err := usecase.GetCategory(ctx, "user-1", "cat-1")

		require.Error(t, err)
		assert.Nil(t, result)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should return error when category not found", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		mockCategoryRepoR.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(nil, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					CategoryRepositoryR: mockCategoryRepoR,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.GetCategory(ctx, "user-1", "cat-1")

		require.Error(t, err)
		assert.Nil(t, result)

		var notFoundErr *domain.EntityNotFoundError
		require.ErrorAs(t, err, &notFoundErr)
	})

	t.Run("should successfully get category", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		category := &domain.Category{
			ID:        "cat-1",
			UserID:    "user-1",
			Name:      "Food",
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		}

		mockCategoryRepoR.EXPECT().
			GetByUserIDAndID(gomock.Any(), "user-1", "cat-1").
			Return(category, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					CategoryRepositoryR: mockCategoryRepoR,
				}
				return fn(ctx, repo)
			})

		result, err := usecase.GetCategory(ctx, "user-1", "cat-1")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "cat-1", result.ID)
		assert.Equal(t, "Food", result.Name)
	})
}

func TestCategoryUseCase_ListCategories(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockIDGen := mocks.NewMockIDGenerator(ctrl)
	clock := clockwork.NewFakeClockAt(fixedTime)
	mockCategoryRepoR := mocks.NewMockCategoryRepositoryR(ctrl)

	usecase := usecases.NewCategoryUseCase(mockDB, mockIDGen, clock)

	t.Run("should return authorization error when user cannot view resource", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-2")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		result, total, err := usecase.ListCategories(ctx, "user-1", nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, total)

		var authErr *domain.AuthorizationError
		require.ErrorAs(t, err, &authErr)
	})

	t.Run("should successfully list categories", func(t *testing.T) {
		ctx := credscontext.WithUserID(context.Background(), "user-1")
		ctx = credscontext.WithUserRole(ctx, domain.UserRoleUser)

		categories := []*domain.Category{
			{
				ID:        "cat-1",
				UserID:    "user-1",
				Name:      "Food",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			{
				ID:        "cat-2",
				UserID:    "user-1",
				Name:      "Transport",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
		}

		mockCategoryRepoR.EXPECT().
			List(gomock.Any(), "user-1", nil, nil).
			Return(categories, 2, nil)

		mockDB.EXPECT().
			DoRead(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(context.Context, interfaces.RepositoriesRead) error) error {
				repo := interfaces.RepositoriesRead{
					CategoryRepositoryR: mockCategoryRepoR,
				}
				return fn(ctx, repo)
			})

		result, total, err := usecase.ListCategories(ctx, "user-1", nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, total)
		assert.Len(t, result, 2)
		assert.Equal(t, "cat-1", result[0].ID)
		assert.Equal(t, "cat-2", result[1].ID)
	})
}
