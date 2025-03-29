package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/jonboulle/clockwork"

	"money_tracker/internal/application/interfaces"
	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
	"money_tracker/internal/pkg/credscontext"
)

type CategoryUseCase struct {
	db          interfaces.Database
	idGenerator interfaces.IDGenerator
	clock       clockwork.Clock
}

func NewCategoryUseCase(
	db interfaces.Database,
	idGenerator interfaces.IDGenerator,
	clock clockwork.Clock,
) *CategoryUseCase {
	return &CategoryUseCase{
		db:          db,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, userID string, input dto.CreateCategoryDTO) (*dto.CategoryDTO, error) {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, &domain.InvalidArgumentError{Argument: "name", Reason: "name is required"}
	}

	now := uc.clock.Now()
	category := &domain.Category{
		ID:        uc.idGenerator.Generate(),
		UserID:    userID,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		existingCategory, err := repo.CategoryRepositoryRW.GetByUserIDAndName(ctx, userID, input.Name)
		if err != nil {
			return fmt.Errorf("check for existing category: %w", err)
		}

		if existingCategory != nil {
			return &domain.EntityAlreadyExistsError{EntityType: "Category", Field: "name", Value: input.Name}
		}

		if err := repo.CategoryRepositoryRW.Create(ctx, category); err != nil {
			return fmt.Errorf("create category: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromCategory(category), nil
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, userID, categoryID string) (*dto.CategoryDTO, error) {
	if !domain.CanViewResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	var category *domain.Category
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		category, err = repo.CategoryRepositoryR.GetByUserIDAndID(ctx, userID, categoryID)
		if err != nil {
			return fmt.Errorf("get category: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, &domain.EntityNotFoundError{EntityType: "Category", ID: categoryID}
	}

	return dto.FromCategory(category), nil
}

func (uc *CategoryUseCase) ListCategories(ctx context.Context, userID string, limit, offset *int) ([]*dto.CategoryDTO, int, error) {
	if !domain.CanViewResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, 0, &domain.AuthorizationError{}
	}

	var categories []*domain.Category
	var total int
	err := uc.db.DoRead(ctx, func(ctx context.Context, repo interfaces.RepositoriesRead) error {
		var err error
		categories, total, err = repo.CategoryRepositoryR.List(ctx, userID, limit, offset)
		if err != nil {
			return fmt.Errorf("list categories: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.CategoryDTO, len(categories))
	for i, cat := range categories {
		result[i] = dto.FromCategory(cat)
	}

	return result, total, nil
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, userID, categoryID string, input dto.UpdateCategoryDTO) (*dto.CategoryDTO, error) {
	if !domain.CanEditResource(credscontext.UserIDFromContext(ctx), credscontext.UserRoleFromContext(ctx), userID) {
		return nil, &domain.AuthorizationError{}
	}

	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return nil, &domain.InvalidArgumentError{Argument: "name", Reason: "name is required"}
	}

	var category *domain.Category
	err := uc.db.DoReadWrite(ctx, func(ctx context.Context, repo interfaces.RepositoriesReadWrite) error {
		var err error
		category, err = repo.CategoryRepositoryRW.GetByUserIDAndID(ctx, userID, categoryID)
		if err != nil {
			return fmt.Errorf("get category: %w", err)
		}

		if category == nil {
			return &domain.EntityNotFoundError{EntityType: "Category", ID: categoryID}
		}

		existingCategory, err := repo.CategoryRepositoryRW.GetByUserIDAndName(ctx, userID, input.Name)
		if err != nil {
			return fmt.Errorf("check for existing category: %w", err)
		}

		if existingCategory != nil && existingCategory.ID != categoryID {
			return &domain.EntityAlreadyExistsError{EntityType: "Category", Field: "name", Value: input.Name}
		}

		category.Name = input.Name
		category.UpdatedAt = uc.clock.Now()

		if err := repo.CategoryRepositoryRW.Update(ctx, category); err != nil {
			return fmt.Errorf("update category: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dto.FromCategory(category), nil
}
