package model

import (
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/uptrace/bun"
)

type Session struct {
	BaseModel

	ID        string `bun:",pk"`
	UserID    string
	Token     string
	ExpiresAt time.Time
	UpdatedAt time.Time
	CreatedAt time.Time

	model[Session, domain.Session, domain.SessionCreateInput, domain.SessionUpdateInput, domain.SessionFilters]
}

func NewSession(db *bun.DB) *Session {
	session := &Session{}
	session.SetDB(db)
	return session
}

func (s Session) id() string {
	return s.ID
}

func (Session) idPrefix() string {
	return "ses_"
}

func (Session) tableName() string {
	return "sessions"
}

func (Session) postprocess(model Session) domain.Session {
	return domain.Session{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		UpdatedAt: model.UpdatedAt,
		CreatedAt: model.CreatedAt,
	}
}
