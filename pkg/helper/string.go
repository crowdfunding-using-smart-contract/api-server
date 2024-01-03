package helper

import "strings"

func ParseString[T any](mapString map[string]T, str string) (T, bool) {
	c, ok := mapString[strings.ToLower(str)]
	return c, ok
}
