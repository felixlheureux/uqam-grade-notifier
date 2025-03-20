package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	UQAMPassword string    `json:"-"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type UserCreateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateInput struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type UserFilters struct {
	Email *string `json:"email,omitempty"`
}
