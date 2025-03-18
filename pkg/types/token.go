package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type SessionClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewTokenClaims(email string, duration time.Duration) *TokenClaims {
	now := time.Now()
	return &TokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}
}

func NewSessionClaims(email string) *SessionClaims {
	now := time.Now()
	return &SessionClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)), // 7 jours
		},
	}
}

func (c *TokenClaims) IsExpired() bool {
	return time.Now().After(c.ExpiresAt.Time)
}
