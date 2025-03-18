package store

import (
	"fmt"
	"strings"

	"github.com/segmentio/ksuid"
)

const (
	CoursePrefix = "crs"
)

// GenerateID génère un ID unique avec un préfixe
func GenerateID(prefix string) string {
	id := ksuid.New().String()
	return fmt.Sprintf("%s_%s", prefix, id)
}

// ParseID extrait le préfixe et l'ID d'une chaîne formatée
func ParseID(id string) (prefix, ksuid string, err error) {
	parts := strings.Split(id, "_")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("format d'ID invalide: %s", id)
	}
	return parts[0], parts[1], nil
}

// ValidateID vérifie si un ID est valide pour un préfixe donné
func ValidateID(id, expectedPrefix string) error {
	prefix, _, err := ParseID(id)
	if err != nil {
		return err
	}
	if prefix != expectedPrefix {
		return fmt.Errorf("préfixe d'ID invalide: attendu %s, reçu %s", expectedPrefix, prefix)
	}
	return nil
}
