package random

import "math/rand"

func NewInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func NewFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
