package dto

import (
	"time"

	usecaseDTO "money_tracker/internal/application/usecases/dto"
)

// CreateUserRequest represents the user registration data
type CreateUserRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

func (r *CreateUserRequest) ToUsecaseDTO() usecaseDTO.CreateUserDTO {
	return usecaseDTO.CreateUserDTO{
		Email:    r.Email,
		Password: r.Password,
	}
}

// UserResponse represents user information
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"user@example.com"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

func NewUserResponseFromUsecaseDTO(dto *usecaseDTO.UserDTO) UserResponse {
	return UserResponse{
		ID:        dto.ID,
		Email:     dto.Email,
		Role:      dto.Role,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
