package domain

import "time"

type Course struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Semester  string    `json:"semester"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type CourseCreateInput struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Semester string `json:"semester"`
}

type CourseUpdateInput struct {
	Code     *string `json:"code,omitempty"`
	Name     *string `json:"name,omitempty"`
	Semester *string `json:"semester,omitempty"`
}

type CourseFilters struct {
	Code     *string `json:"code,omitempty"`
	Name     *string `json:"name,omitempty"`
	Semester *string `json:"semester,omitempty"`
}
