package dto

import (
	"time"

	usecaseDTO "money_tracker/internal/application/usecases/dto"
)

// CreateCategoryRequest represents the data for creating a category
type CreateCategoryRequest struct {
	Name string `json:"name" example:"Groceries"`
}

func (r *CreateCategoryRequest) ToUsecaseDTO() usecaseDTO.CreateCategoryDTO {
	return usecaseDTO.CreateCategoryDTO{
		Name: r.Name,
	}
}

// UpdateCategoryRequest represents the data for updating a category
type UpdateCategoryRequest struct {
	Name string `json:"name" example:"Food & Groceries"`
}

func (r *UpdateCategoryRequest) ToUsecaseDTO() usecaseDTO.UpdateCategoryDTO {
	return usecaseDTO.UpdateCategoryDTO{
		Name: r.Name,
	}
}

// CategoryListResponse represents a paginated list of categories
type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int                `json:"total" example:"10"`
	Limit      *int               `json:"limit,omitempty" example:"20"`
	Offset     *int               `json:"offset,omitempty" example:"0"`
}

// CategoryResponse represents category information
type CategoryResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name      string    `json:"name" example:"Groceries"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func NewCategoryResponseFromUsecaseDTO(dto *usecaseDTO.CategoryDTO) CategoryResponse {
	return CategoryResponse{
		ID:        dto.ID,
		UserID:    dto.UserID,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
