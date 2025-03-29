package dto_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money_tracker/internal/application/usecases/dto"
	httpDTO "money_tracker/internal/delivery/http/dto"
)

func TestLoginRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.LoginRequest
		expected dto.LoginDTO
	}{
		{
			name: "should correctly convert login request",
			input: &httpDTO.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expected: dto.LoginDTO{
				Email:    "test@example.com",
				Password: "password123",
			},
		},
		{
			name: "should correctly convert login request with empty fields",
			input: &httpDTO.LoginRequest{
				Email:    "",
				Password: "",
			},
			expected: dto.LoginDTO{
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

func TestNewLoginResponseFromUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *dto.AuthResponseDTO
		expected httpDTO.LoginResponse
	}{
		{
			name: "should correctly convert auth response",
			input: &dto.AuthResponseDTO{
				Token: "jwt_token_123",
				User: &dto.UserDTO{
					ID:        "user-123",
					Email:     "test@example.com",
					Role:      "user",
					CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				},
			},
			expected: httpDTO.LoginResponse{
				Token: "jwt_token_123",
				User: httpDTO.UserResponse{
					ID:        "user-123",
					Email:     "test@example.com",
					Role:      "user",
					CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "should correctly convert admin auth response",
			input: &dto.AuthResponseDTO{
				Token: "admin_token_456",
				User: &dto.UserDTO{
					ID:        "admin-456",
					Email:     "admin@example.com",
					Role:      "admin",
					CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
				},
			},
			expected: httpDTO.LoginResponse{
				Token: "admin_token_456",
				User: httpDTO.UserResponse{
					ID:        "admin-456",
					Email:     "admin@example.com",
					Role:      "admin",
					CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := httpDTO.NewLoginResponseFromUsecaseDTO(tc.input)

			require.NotNil(t, result)
			assert.Equal(t, tc.expected.Token, result.Token)
			assert.Equal(t, tc.expected.User.ID, result.User.ID)
			assert.Equal(t, tc.expected.User.Email, result.User.Email)
			assert.Equal(t, tc.expected.User.Role, result.User.Role)
			assert.Equal(t, tc.expected.User.CreatedAt, result.User.CreatedAt)
			assert.Equal(t, tc.expected.User.UpdatedAt, result.User.UpdatedAt)
		})
	}
}
