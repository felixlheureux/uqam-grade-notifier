package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func New(connString string) (*DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture de la connexion à la base de données: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification de la connexion à la base de données: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
