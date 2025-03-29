package password_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"money_tracker/internal/infrastructure/password"
)

func TestBcryptHasher_Verify(t *testing.T) {
	hasher := password.NewBcryptHasher()

	testCases := []struct {
		name     string
		password string
	}{
		{
			name:     "should verify correct password",
			password: "password123",
		},
		{
			name:     "should verify complex password",
			password: "P@ssw0rd!2024#Complex",
		},
		{
			name:     "should verify short password",
			password: "abc",
		},
		{
			name:     "should verify password with unicode",
			password: "пароль123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := hasher.Hash(tc.password)
			require.NoError(t, err)

			err = hasher.Verify(tc.password, hash)

			require.NoError(t, err)
		})
	}
}

func TestBcryptHasher_Verify_WrongPassword(t *testing.T) {
	hasher := password.NewBcryptHasher()

	testCases := []struct {
		name            string
		correctPassword string
		wrongPassword   string
	}{
		{
			name:            "should fail with completely different password",
			correctPassword: "password123",
			wrongPassword:   "wrongpassword",
		},
		{
			name:            "should fail with similar password",
			correctPassword: "password123",
			wrongPassword:   "password124",
		},
		{
			name:            "should fail with case difference",
			correctPassword: "password123",
			wrongPassword:   "Password123",
		},
		{
			name:            "should fail with empty password",
			correctPassword: "password123",
			wrongPassword:   "",
		},
		{
			name:            "should fail with extra characters",
			correctPassword: "password123",
			wrongPassword:   "password123 ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := hasher.Hash(tc.correctPassword)
			require.NoError(t, err)

			err = hasher.Verify(tc.correctPassword, hash)
			require.NoError(t, err)

			err = hasher.Verify(tc.wrongPassword, hash)

			require.Error(t, err)
		})
	}
}
