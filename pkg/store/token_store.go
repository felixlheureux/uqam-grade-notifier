package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
)

type TokenStore struct {
	db *db.DB
}

func NewTokenStore(db *db.DB) (*TokenStore, error) {
	if err := createTokenTables(db); err != nil {
		return nil, err
	}
	return &TokenStore{db: db}, nil
}

func createTokenTables(db *db.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS tokens (
			id VARCHAR(50) PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			token VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(email, token)
		);
	`

	_, err := db.Exec(query)
	return err
}

func (s *TokenStore) SaveToken(email, token string) error {
	id := GenerateID("tok")
	expiresAt := time.Now().Add(15 * time.Minute)

	_, err := squirrel.Insert("tokens").
		Columns("id", "email", "token", "expires_at").
		Values(id, email, token, expiresAt).
		RunWith(s.db).
		Exec()
	if err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde du token: %w", err)
	}
	return nil
}

func (s *TokenStore) GetToken(email string) (string, error) {
	var token string
	var expiresAt time.Time
	err := squirrel.Select("token", "expires_at").
		From("tokens").
		Where(squirrel.Eq{"email": email}).
		OrderBy("created_at DESC").
		Limit(1).
		RunWith(s.db).
		QueryRow().
		Scan(&token, &expiresAt)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération du token: %w", err)
	}

	if time.Now().After(expiresAt) {
		return "", nil
	}

	return token, nil
}

func (s *TokenStore) DeleteToken(email string) error {
	result, err := squirrel.Delete("tokens").
		Where(squirrel.Eq{"email": email}).
		RunWith(s.db).
		Exec()
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("token non trouvé")
	}

	return nil
}

func (s *TokenStore) SaveSessionToken(email, token string) error {
	id := GenerateID("ses")
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err := squirrel.Insert("tokens").
		Columns("id", "email", "token", "expires_at").
		Values(id, email, token, expiresAt).
		RunWith(s.db).
		Exec()
	if err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde du token de session: %w", err)
	}
	return nil
}

func (s *TokenStore) ValidateSessionToken(token string) (string, error) {
	var email string
	var expiresAt time.Time
	err := squirrel.Select("email", "expires_at").
		From("tokens").
		Where(squirrel.Eq{"token": token}).
		RunWith(s.db).
		QueryRow().
		Scan(&email, &expiresAt)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("token de session non trouvé")
	}
	if err != nil {
		return "", fmt.Errorf("erreur lors de la validation du token de session: %w", err)
	}

	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("token de session expiré")
	}

	return email, nil
}
