package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"money_tracker/internal/domain"
)

func TestCanViewResource(t *testing.T) {
	testCases := []struct {
		name            string
		userID          string
		role            domain.UserRole
		resourceOwnerID string
		expected        bool
	}{
		{
			name:            "should allow user to view their own resource",
			userID:          "user-123",
			role:            domain.UserRoleUser,
			resourceOwnerID: "user-123",
			expected:        true,
		},
		{
			name:            "should not allow user to view other user's resource",
			userID:          "user-123",
			role:            domain.UserRoleUser,
			resourceOwnerID: "user-456",
			expected:        false,
		},
		{
			name:            "should allow admin to view their own resource",
			userID:          "admin-123",
			role:            domain.UserRoleAdmin,
			resourceOwnerID: "admin-123",
			expected:        true,
		},
		{
			name:            "should allow admin to view other user's resource",
			userID:          "admin-123",
			role:            domain.UserRoleAdmin,
			resourceOwnerID: "user-456",
			expected:        true,
		},
		{
			name:            "should not allow empty user to view resource",
			userID:          "",
			role:            domain.UserRoleUser,
			resourceOwnerID: "user-123",
			expected:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.CanViewResource(tc.userID, tc.role, tc.resourceOwnerID)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCanEditResource(t *testing.T) {
	testCases := []struct {
		name            string
		userID          string
		role            domain.UserRole
		resourceOwnerID string
		expected        bool
	}{
		{
			name:            "should allow user to edit their own resource",
			userID:          "user-123",
			role:            domain.UserRoleUser,
			resourceOwnerID: "user-123",
			expected:        true,
		},
		{
			name:            "should not allow user to edit other user's resource",
			userID:          "user-123",
			role:            domain.UserRoleUser,
			resourceOwnerID: "user-456",
			expected:        false,
		},
		{
			name:            "should allow admin to edit their own resource",
			userID:          "admin-123",
			role:            domain.UserRoleAdmin,
			resourceOwnerID: "admin-123",
			expected:        true,
		},
		{
			name:            "should not allow admin to edit other user's resource",
			userID:          "admin-123",
			role:            domain.UserRoleAdmin,
			resourceOwnerID: "user-456",
			expected:        false,
		},
		{
			name:            "should not allow empty user to edit resource",
			userID:          "",
			role:            domain.UserRoleUser,
			resourceOwnerID: "user-123",
			expected:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.CanEditResource(tc.userID, tc.role, tc.resourceOwnerID)

			assert.Equal(t, tc.expected, result)
		})
	}
}
