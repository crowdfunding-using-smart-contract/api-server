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

func TestRandomFloat32Number(t *testing.T) {
	t.Run("Test NewFloat32", func(t *testing.T) {
		min := float32(1.0)
		max := float32(10.0)
		number := NewFloat32(min, max)
		require.GreaterOrEqual(t, number, min)
		require.LessOrEqual(t, number, max)
	})
}
