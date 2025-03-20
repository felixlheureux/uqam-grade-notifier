package service

import (
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/jwt"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/model"
)

type AuthService struct {
	tokenModel   *model.Token
	tokenManager *jwt.Manager
}

func NewAuthService(tokenModel *model.Token, tokenManager *jwt.Manager) *AuthService {
	return &AuthService{
		tokenModel:   tokenModel,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) GenerateToken(email string, duration time.Duration) (string, error) {
	token, err := s.tokenManager.GenerateToken(email, duration)
	if err != nil {
		return "", domain.ErrTokenGenerationFailed(err)
	}
	return token, nil
}

func (s *AuthService) ValidateToken(token string) (*domain.TokenClaims, error) {
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return nil, domain.ErrTokenValidationFailed(err)
	}
	return claims, nil
}

func (s *AuthService) SaveToken(email, token string) error {
	if err := s.tokenModel.SaveToken(email, token); err != nil {
		return domain.ErrTokenStorageFailed(err)
	}
	return nil
}

func (s *AuthService) GetToken(email string) (string, error) {
	token, err := s.tokenModel.GetToken(email)
	if err != nil {
		return "", domain.ErrTokenStorageFailed(err)
	}
	return token, nil
}

func (s *AuthService) DeleteToken(email string) error {
	if err := s.tokenModel.DeleteToken(email); err != nil {
		return domain.ErrTokenStorageFailed(err)
	}
	return nil
}

func (s *AuthService) GenerateSessionToken(email string) (string, error) {
	token, err := s.tokenManager.GenerateSessionToken(email)
	if err != nil {
		return "", domain.ErrTokenGenerationFailed(err)
	}
	return token, nil
}

func (s *AuthService) ValidateSessionToken(token string) (string, error) {
	email, err := s.tokenModel.ValidateSessionToken(token)
	if err != nil {
		return "", domain.ErrTokenValidationFailed(err)
	}
	return email, nil
}
