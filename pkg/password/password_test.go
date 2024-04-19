package password

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	t.Run("Test PasswordCondition", func(t *testing.T) {
		t.Run("Test PasswordCondition with valid password", func(t *testing.T) {
			password := "Password1!"
			require.True(t, PasswordCondition(password))
		})

		t.Run("Test PasswordCondition with invalid password", func(t *testing.T) {
			password := "password"
			require.False(t, PasswordCondition(password))
		})
	})

	t.Run("Test HashPassword", func(t *testing.T) {
		t.Run("Test HashPassword with valid password", func(t *testing.T) {
			password := "Password1!"
			hashedPassword, err := HashPassword(password)
			require.NoError(t, err)
			require.NotEmpty(t, hashedPassword)
		})

		t.Run("Test HashPassword with invalid password", func(t *testing.T) {
			password := "password"
			hashedPassword, err := HashPassword(password)
			require.Error(t, err)
			require.Empty(t, hashedPassword)
		})
	})

	t.Run("Test CheckPassword", func(t *testing.T) {
		t.Run("Test CheckPassword with valid password", func(t *testing.T) {
			password := "Password1!"
			hashedPassword, err := HashPassword(password)
			require.NoError(t, err)
			require.NotEmpty(t, hashedPassword)

			err = CheckPassword(password, hashedPassword)
			require.NoError(t, err)
		})

		t.Run("Test CheckPassword with invalid password", func(t *testing.T) {
			password := "Password1!"
			hashedPassword, err := HashPassword(password)
			require.NoError(t, err)
			require.NotEmpty(t, hashedPassword)

			wrongPassword := "password"
			err = CheckPassword(wrongPassword, hashedPassword)
			require.Error(t, err)
		})
	})
}
