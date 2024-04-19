package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	t.Run("Test ParseString", func(t *testing.T) {
		mapString := map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
		str := "one"
		got, ok := ParseString(mapString, str)
		require.Equal(t, ok, true)
		require.Equal(t, got, 1)
	})
}
