package service

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/model"
)

type SessionService struct {
	sessionModel *model.Session
}

func NewSessionService(sessionModel *model.Session) *SessionService {
	return &SessionService{
		sessionModel: sessionModel,
	}
}

func (s *SessionService) Get(ctx domain.Context) ([]domain.Session, error) {
	sessions, err := s.sessionModel.Get("", ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUnexpected(err)
	}
	return []domain.Session{sessions}, nil
}

func (s *SessionService) FindOne(ctx domain.Context, id string) (*domain.Session, error) {
	session, err := s.sessionModel.Get(id, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrNotFound(err)
	}
	return &session, nil
}

func (s *SessionService) Create(ctx domain.Context, input domain.SessionCreateInput) (*domain.Session, error) {
	session, err := s.sessionModel.Create(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUnexpected(err)
	}
	return &session, nil
}

func (s *SessionService) Update(ctx domain.Context, id string, input domain.SessionUpdateInput) (*domain.Session, error) {
	session, err := s.sessionModel.Update(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUnexpected(err)
	}
	return &session, nil
}

func (s *SessionService) Delete(ctx domain.Context, id string) error {
	_, err := s.sessionModel.Delete(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrUnexpected(err)
	}
	return nil
}

func (s *SessionService) Destroy(ctx domain.Context, id string) error {
	_, err := s.sessionModel.Destroy(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrUnexpected(err)
	}
	return nil
}
