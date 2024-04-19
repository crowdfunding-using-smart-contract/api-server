package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	t.Run("Test NewString", func(t *testing.T) {
		length := 10
		str := NewString(length)
		require.Equal(t, length, len(str))
	})
}
