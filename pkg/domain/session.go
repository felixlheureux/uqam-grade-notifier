package domain

import "time"

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type SessionCreateInput struct {
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SessionUpdateInput struct {
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type SessionFilters struct {
	UserID *string `json:"userId,omitempty"`
	Token  *string `json:"token,omitempty"`
}
