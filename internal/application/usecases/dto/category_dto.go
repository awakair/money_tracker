package dto

import (
	"time"

	"money_tracker/internal/domain"
)

type CreateCategoryDTO struct {
	Name string
}

type UpdateCategoryDTO struct {
	Name string
}

type CategoryDTO struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromCategory(category *domain.Category) *CategoryDTO {
	return &CategoryDTO{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
