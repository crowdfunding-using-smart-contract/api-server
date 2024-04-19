package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomEmail(t *testing.T) {
	t.Run("Test NewEmail", func(t *testing.T) {
		email := NewEmail()
		require.NotEmpty(t, email)
		require.Contains(t, email, "@")
	})
}
