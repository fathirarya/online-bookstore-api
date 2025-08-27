package auth

import "time"

type JWTConfig struct {
	SecretKey      string
	Issuer         string
	ExpireDuration time.Duration
}
