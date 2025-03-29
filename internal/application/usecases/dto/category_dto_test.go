package dto_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money_tracker/internal/application/usecases/dto"
	"money_tracker/internal/domain"
)

func TestFromCategory(t *testing.T) {
	testCases := []struct {
		name     string
		input    *domain.Category
		expected *dto.CategoryDTO
	}{
		{
			name: "should correctly convert category",
			input: &domain.Category{
				ID:        "cat-123",
				UserID:    "user-456",
				Name:      "Food",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			expected: &dto.CategoryDTO{
				ID:        "cat-123",
				UserID:    "user-456",
				Name:      "Food",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "should correctly convert category with different timestamps",
			input: &domain.Category{
				ID:        "cat-789",
				UserID:    "user-123",
				Name:      "Transport",
				CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
			},
			expected: &dto.CategoryDTO{
				ID:        "cat-789",
				UserID:    "user-123",
				Name:      "Transport",
				CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := dto.FromCategory(tc.input)

			require.NotNil(t, result)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.UserID, result.UserID)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
