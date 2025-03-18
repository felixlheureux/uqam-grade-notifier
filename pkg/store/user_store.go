package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/Masterminds/squirrel"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
)

type UserStore struct {
	db    *db.DB
	key   []byte
	block cipher.Block
}

func NewUserStore(db *db.DB, key []byte) (*UserStore, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du chiffrement: %w", err)
	}

	if err := createUserTables(db); err != nil {
		return nil, err
	}

	return &UserStore{
		db:    db,
		key:   key,
		block: block,
	}, nil
}

func createUserTables(db *db.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(50) PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			uqam_username VARCHAR(255) NOT NULL,
			uqam_password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	return err
}

func (s *UserStore) encrypt(plaintext []byte) (string, error) {
	gcm, err := cipher.NewGCM(s.block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *UserStore) decrypt(encrypted string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(s.block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext trop court")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (s *UserStore) SaveUser(email, uqamUsername, uqamPassword string) error {
	encryptedPassword, err := s.encrypt([]byte(uqamPassword))
	if err != nil {
		return fmt.Errorf("erreur lors du chiffrement du mot de passe: %w", err)
	}

	id := GenerateID("usr")
	_, err = squirrel.Insert("users").
		Columns("id", "email", "uqam_username", "uqam_password").
		Values(id, email, uqamUsername, encryptedPassword).
		RunWith(s.db).
		Exec()
	if err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde de l'utilisateur: %w", err)
	}

	return nil
}

func (s *UserStore) GetUser(email string) (string, string, error) {
	var uqamUsername, encryptedPassword string
	err := squirrel.Select("uqam_username", "uqam_password").
		From("users").
		Where(squirrel.Eq{"email": email}).
		RunWith(s.db).
		QueryRow().
		Scan(&uqamUsername, &encryptedPassword)
	if err == sql.ErrNoRows {
		return "", "", fmt.Errorf("utilisateur non trouvé")
	}
	if err != nil {
		return "", "", fmt.Errorf("erreur lors de la récupération de l'utilisateur: %w", err)
	}

	uqamPassword, err := s.decrypt(encryptedPassword)
	if err != nil {
		return "", "", fmt.Errorf("erreur lors du déchiffrement du mot de passe: %w", err)
	}

	return uqamUsername, string(uqamPassword), nil
}

func (s *UserStore) GetAllUsers() ([]string, error) {
	rows, err := squirrel.Select("email").
		From("users").
		RunWith(s.db).
		Query()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des utilisateurs: %w", err)
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture des emails: %w", err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}
