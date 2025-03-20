package model

import (
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/uptrace/bun"
)

type Grade struct {
	BaseModel

	ID        string `bun:",pk"`
	UserID    string
	Semester  string
	Course    string
	Grade     string
	UpdatedAt time.Time
	CreatedAt time.Time

	model[Grade, domain.Grade, domain.GradeCreateInput, domain.GradeUpdateInput, domain.GradeFilters]
}

func NewGrade(db *bun.DB) *Grade {
	grade := &Grade{}
	grade.SetDB(db)
	return grade
}

func (g Grade) id() string {
	return g.ID
}

func (Grade) idPrefix() string {
	return "grd_"
}

func (Grade) tableName() string {
	return "grades"
}

func (Grade) postprocess(model Grade) domain.Grade {
	return domain.Grade{
		ID:        model.ID,
		UserID:    model.UserID,
		Semester:  model.Semester,
		Course:    model.Course,
		Grade:     model.Grade,
		UpdatedAt: model.UpdatedAt,
		CreatedAt: model.CreatedAt,
	}
}
