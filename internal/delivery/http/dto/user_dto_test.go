package dto_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"money_tracker/internal/application/usecases/dto"
	httpDTO "money_tracker/internal/delivery/http/dto"
)

func TestCreateUserRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.CreateUserRequest
		expected dto.CreateUserDTO
	}{
		{
			name: "should correctly convert create user request",
			input: &httpDTO.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expected: dto.CreateUserDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
		},
		{
			name: "should correctly convert create user request with empty fields",
			input: &httpDTO.CreateUserRequest{
				Email:    "",
				Password: "",
			},
			expected: dto.CreateUserDTO{
				Email:    "",
				Password: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.Email, result.Email)
			assert.Equal(t, tc.expected.Password, result.Password)
		})
	}
}

func TestNewUserResponseFromUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *dto.UserDTO
		expected httpDTO.UserResponse
	}{
		{
			name: "should correctly convert user DTO to response",
			input: &dto.UserDTO{
				ID:        "user-123",
				Email:     "test@example.com",
				Role:      "user",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			expected: httpDTO.UserResponse{
				ID:        "user-123",
				Email:     "test@example.com",
				Role:      "user",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "should correctly convert admin user DTO to response",
			input: &dto.UserDTO{
				ID:        "admin-456",
				Email:     "admin@example.com",
				Role:      "admin",
				CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
			},
			expected: httpDTO.UserResponse{
				ID:        "admin-456",
				Email:     "admin@example.com",
				Role:      "admin",
				CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := httpDTO.NewUserResponseFromUsecaseDTO(tc.input)

			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Email, result.Email)
			assert.Equal(t, tc.expected.Role, result.Role)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
