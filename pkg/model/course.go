package model

import (
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/uptrace/bun"
)

type Course struct {
	BaseModel

	ID        string `bun:",pk"`
	Code      string
	Name      string
	Semester  string
	UpdatedAt time.Time
	CreatedAt time.Time

	model[Course, domain.Course, domain.CourseCreateInput, domain.CourseUpdateInput, domain.CourseFilters]
}

func NewCourse(db *bun.DB) *Course {
	course := &Course{}
	course.SetDB(db)
	return course
}

func (c Course) id() string {
	return c.ID
}

func (Course) idPrefix() string {
	return "crs_"
}

func (Course) tableName() string {
	return "courses"
}

func (Course) postprocess(model Course) domain.Course {
	return domain.Course{
		ID:        model.ID,
		Code:      model.Code,
		Name:      model.Name,
		Semester:  model.Semester,
		UpdatedAt: model.UpdatedAt,
		CreatedAt: model.CreatedAt,
	}
}
