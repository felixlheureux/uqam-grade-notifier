package domain

import "time"

type Grade struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Semester  string    `json:"semester"`
	Course    string    `json:"course"`
	Grade     string    `json:"grade"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type GradeCreateInput struct {
	UserID   string `json:"userId"`
	Semester string `json:"semester"`
	Course   string `json:"course"`
	Grade    string `json:"grade"`
}

type GradeUpdateInput struct {
	Grade string `json:"grade"`
}

type GradeFilters struct {
	UserID   *string `json:"userId,omitempty"`
	Semester *string `json:"semester,omitempty"`
	Course   *string `json:"course,omitempty"`
}
