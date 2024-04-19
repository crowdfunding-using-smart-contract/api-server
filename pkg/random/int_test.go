package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomIntegerNumber(t *testing.T) {
	t.Run("Test NewInt", func(t *testing.T) {
		min := 1
		max := 10
		number := NewInt(min, max)
		require.GreaterOrEqual(t, number, min)
		require.LessOrEqual(t, number, max)
	})
}
