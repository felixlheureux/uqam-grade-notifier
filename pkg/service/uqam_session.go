package service

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/model"
)

type UQAMSessionService struct {
	uqamSessionModel *model.UQAMSession
}

func NewUQAMSessionService(uqamSessionModel *model.UQAMSession) *UQAMSessionService {
	return &UQAMSessionService{
		uqamSessionModel: uqamSessionModel,
	}
}

func (s *UQAMSessionService) Get(ctx domain.Context) ([]domain.UQAMSession, error) {
	sessions, err := s.uqamSessionModel.Get("", ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUnexpected(err)
	}
	return []domain.UQAMSession{sessions}, nil
}

func (s *UQAMSessionService) FindOne(ctx domain.Context, id string) (*domain.UQAMSession, error) {
	session, err := s.uqamSessionModel.Get(id, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrNotFound(err)
	}
	return &session, nil
}

func (s *UQAMSessionService) Create(ctx domain.Context, input domain.UQAMSessionCreateInput) (*domain.UQAMSession, error) {
	session, err := s.uqamSessionModel.Create(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUnexpected(err)
	}
	return &session, nil
}

func (s *UQAMSessionService) Update(ctx domain.Context, id string, input domain.UQAMSessionUpdateInput) (*domain.UQAMSession, error) {
	session, err := s.uqamSessionModel.Update(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUnexpected(err)
	}
	return &session, nil
}

func (s *UQAMSessionService) Delete(ctx domain.Context, id string) error {
	_, err := s.uqamSessionModel.Delete(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrUnexpected(err)
	}
	return nil
}

func (s *UQAMSessionService) Destroy(ctx domain.Context, id string) error {
	_, err := s.uqamSessionModel.Destroy(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrUnexpected(err)
	}
	return nil
}
