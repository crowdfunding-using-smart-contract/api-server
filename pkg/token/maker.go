package token

import "time"

type Maker interface {
	CreateToken(userId uint, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
