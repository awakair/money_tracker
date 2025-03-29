package dto_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
)

func TestFromUser(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domain.User
		expected *dto.UserDTO
	}{
		{
			name: "should correctly convert regular user",
			input: &domain.User{
				ID:           "user-123",
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
				Role:         domain.UserRoleUser,
				CreatedAt:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			expected: &dto.UserDTO{
				ID:        "user-123",
				Email:     "test@example.com",
				Role:      "user",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "should correctly convert admin user",
			input: &domain.User{
				ID:           "admin-456",
				Email:        "admin@example.com",
				PasswordHash: "hashed_password",
				Role:         domain.UserRoleAdmin,
				CreatedAt:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt:    time.Date(2025, 1, 16, 12, 0, 0, 0, time.UTC),
			},
			expected: &dto.UserDTO{
				ID:        "admin-456",
				Email:     "admin@example.com",
				Role:      "admin",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 16, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := dto.FromUser(tc.input)

			require.NotNil(t, result)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Email, result.Email)
			assert.Equal(t, tc.expected.Role, result.Role)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
