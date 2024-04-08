package random

import "math/rand"

func NewInt(min, max int) int {
	return min + rand.Intn(max-min)
}
