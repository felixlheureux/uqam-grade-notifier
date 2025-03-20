package domain

import "time"

// TestEntity représente une entité de test
type TestEntity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TestCreateInput représente les données d'entrée pour la création d'une entité de test
type TestCreateInput struct {
	Name string `json:"name"`
}

// TestUpdateInput représente les données d'entrée pour la mise à jour d'une entité de test
type TestUpdateInput struct {
	Name string `json:"name"`
}

// TestFilters représente les filtres pour la recherche d'entités de test
type TestFilters struct {
	Name string `json:"name"`
}
