package jwt

import (
	"fmt"
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey []byte
}

func NewManager(secretKey string) *Manager {
	return &Manager{
		secretKey: []byte(secretKey),
	}
}

func (m *Manager) GenerateToken(email string, duration time.Duration) (string, error) {
	claims := domain.NewTokenClaims(email, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *Manager) GenerateSessionToken(email string) (string, error) {
	claims := domain.NewSessionClaims(email)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

func (m *Manager) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (m *Manager) ValidateSessionToken(tokenString string) (*domain.SessionClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.SessionClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.SessionClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid session token")
}
