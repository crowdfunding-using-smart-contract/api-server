package random

import "fmt"

func NewEmail() string {
	return fmt.Sprintf("%s@%s.com", NewString(10), NewString(5))
}
