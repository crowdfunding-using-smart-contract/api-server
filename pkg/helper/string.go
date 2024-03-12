package helper

import (
	"math/rand"
	"strings"
)

func ParseString[T any](mapString map[string]T, str string) (T, bool) {
	c, ok := mapString[strings.ToLower(str)]
	return c, ok
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
