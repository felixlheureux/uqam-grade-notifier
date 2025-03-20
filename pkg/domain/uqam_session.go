package domain

import "time"

type UQAMSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type UQAMSessionCreateInput struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

type UQAMSessionUpdateInput struct {
	Token *string `json:"token,omitempty"`
}

type UQAMSessionFilters struct {
	UserID *string `json:"userId,omitempty"`
}
