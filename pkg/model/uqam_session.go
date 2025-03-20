package model

import (
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/uptrace/bun"
)

type UQAMSession struct {
	BaseModel

	ID        string `bun:",pk"`
	UserID    string
	Token     string
	UpdatedAt time.Time
	CreatedAt time.Time

	model[UQAMSession, domain.UQAMSession, domain.UQAMSessionCreateInput, domain.UQAMSessionUpdateInput, domain.UQAMSessionFilters]
}

func NewUQAMSession(db *bun.DB) *UQAMSession {
	session := &UQAMSession{}
	session.SetDB(db)
	return session
}

func (s UQAMSession) id() string {
	return s.ID
}

func (UQAMSession) idPrefix() string {
	return "uqam_"
}

func (UQAMSession) tableName() string {
	return "uqam_sessions"
}

func (UQAMSession) postprocess(model UQAMSession) domain.UQAMSession {
	return domain.UQAMSession{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		UpdatedAt: model.UpdatedAt,
		CreatedAt: model.CreatedAt,
	}
}
