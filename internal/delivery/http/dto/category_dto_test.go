package dto_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"money_tracker/internal/application/usecases/dto"
	httpDTO "money_tracker/internal/delivery/http/dto"
)

func TestCreateCategoryRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.CreateCategoryRequest
		expected dto.CreateCategoryDTO
	}{
		{
			name: "should correctly convert create category request",
			input: &httpDTO.CreateCategoryRequest{
				Name: "Food",
			},
			expected: dto.CreateCategoryDTO{
				Name: "Food",
			},
		},
		{
			name: "should correctly convert create category request with empty name",
			input: &httpDTO.CreateCategoryRequest{
				Name: "",
			},
			expected: dto.CreateCategoryDTO{
				Name: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.Name, result.Name)
		})
	}
}

func TestUpdateCategoryRequest_ToUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *httpDTO.UpdateCategoryRequest
		expected dto.UpdateCategoryDTO
	}{
		{
			name: "should correctly convert update category request",
			input: &httpDTO.UpdateCategoryRequest{
				Name: "Updated Food",
			},
			expected: dto.UpdateCategoryDTO{
				Name: "Updated Food",
			},
		},
		{
			name: "should correctly convert update category request with empty name",
			input: &httpDTO.UpdateCategoryRequest{
				Name: "",
			},
			expected: dto.UpdateCategoryDTO{
				Name: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToUsecaseDTO()

			assert.Equal(t, tc.expected.Name, result.Name)
		})
	}
}

func TestNewCategoryResponseFromUsecaseDTO(t *testing.T) {
	testCases := []struct {
		name     string
		input    *dto.CategoryDTO
		expected httpDTO.CategoryResponse
	}{
		{
			name: "should correctly convert category DTO to response",
			input: &dto.CategoryDTO{
				ID:        "cat-123",
				UserID:    "user-456",
				Name:      "Food",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
			expected: httpDTO.CategoryResponse{
				ID:        "cat-123",
				UserID:    "user-456",
				Name:      "Food",
				CreatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "should correctly convert category DTO with different timestamps",
			input: &dto.CategoryDTO{
				ID:        "cat-789",
				UserID:    "user-123",
				Name:      "Transport",
				CreatedAt: time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2025, 1, 20, 15, 30, 0, 0, time.UTC),
			},
			expected: httpDTO.CategoryResponse{
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
			result := httpDTO.NewCategoryResponseFromUsecaseDTO(tc.input)

			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.UserID, result.UserID)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tc.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}
