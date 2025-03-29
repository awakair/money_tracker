package dto

import (
	"time"

	"money_tracker/internal/domain"
)

type CreateUserDTO struct {
	Email    string
	Password string
}

type UserDTO struct {
	ID        string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromUser(user *domain.User) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
