package token

import "time"

type TokenGenerator interface {
	Generate() ([]byte, error)
	TTL() time.Duration
}

type TokenRefresher interface {
	Token() ([]byte, bool, error)
	Stop()
}
