package jwt_test

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"money_tracker/internal/domain"
	"money_tracker/internal/infrastructure/jwt"
)

func TestJWTService_Generate(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	clock := clockwork.NewFakeClockAt(fixedTime)

	testCases := []struct {
		name       string
		secret     string
		expiration time.Duration
		userID     string
		role       domain.UserRole
		wantError  bool
	}{
		{
			name:       "should generate token for regular user",
			secret:     "test-secret-key",
			expiration: 24 * time.Hour,
			userID:     "user-123",
			role:       domain.UserRoleUser,
			wantError:  false,
		},
		{
			name:       "should generate token for admin user",
			secret:     "test-secret-key",
			expiration: 24 * time.Hour,
			userID:     "admin-456",
			role:       domain.UserRoleAdmin,
			wantError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := jwt.NewJWTService(tc.secret, tc.expiration, clock)

			token, err := service.Generate(tc.userID, tc.role)

			if tc.wantError {
				require.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	fixedTime := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	clock := clockwork.NewFakeClockAt(fixedTime)
	secret := "test-secret-key"
	expiration := 24 * time.Hour
	service := jwt.NewJWTService(secret, expiration, clock)

	t.Run("should validate valid token for regular user", func(t *testing.T) {
		token, err := service.Generate("user-123", domain.UserRoleUser)
		require.NoError(t, err)

		userID, role, err := service.ValidateToken(token)

		require.NoError(t, err)
		assert.Equal(t, "user-123", userID)
		assert.Equal(t, domain.UserRoleUser, role)
	})

	t.Run("should validate valid token for admin user", func(t *testing.T) {
		token, err := service.Generate("admin-456", domain.UserRoleAdmin)
		require.NoError(t, err)

		userID, role, err := service.ValidateToken(token)

		require.NoError(t, err)
		assert.Equal(t, "admin-456", userID)
		assert.Equal(t, domain.UserRoleAdmin, role)
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		userID, role, err := service.ValidateToken("invalid-token")

		require.Error(t, err)
		assert.Empty(t, userID)
		assert.Empty(t, role)
	})

	t.Run("should return error for empty token", func(t *testing.T) {
		userID, role, err := service.ValidateToken("")

		require.Error(t, err)
		assert.Empty(t, userID)
		assert.Empty(t, role)
	})

	t.Run("should return error for token with wrong secret", func(t *testing.T) {
		wrongService := jwt.NewJWTService("wrong-secret", expiration, clock)
		token, err := service.Generate("user-789", domain.UserRoleUser)
		require.NoError(t, err)

		userID, role, err := wrongService.ValidateToken(token)

		require.Error(t, err)
		assert.Empty(t, userID)
		assert.Empty(t, role)
	})

	t.Run("should return error for expired token", func(t *testing.T) {
		expiredClock := clockwork.NewFakeClockAt(fixedTime)
		shortService := jwt.NewJWTService(secret, 1*time.Hour, expiredClock)

		token, err := shortService.Generate("user-999", domain.UserRoleUser)
		require.NoError(t, err)

		expiredClock.Advance(2 * time.Hour)

		userID, role, err := shortService.ValidateToken(token)

		require.Error(t, err)
		assert.Empty(t, userID)
		assert.Empty(t, role)
	})
}
