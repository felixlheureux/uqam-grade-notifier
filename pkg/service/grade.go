package service

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/model"
)

type GradeService struct {
	gradeModel *model.Grade
}

func NewGradeService(gradeModel *model.Grade) *GradeService {
	return &GradeService{
		gradeModel: gradeModel,
	}
}

func (s *GradeService) Get(ctx domain.Context) ([]domain.Grade, error) {
	grades, err := s.gradeModel.Get("", ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrGradeValidationFailed(err)
	}
	return []domain.Grade{grades}, nil
}

func (s *GradeService) FindOne(ctx domain.Context, id string) (*domain.Grade, error) {
	grade, err := s.gradeModel.Get(id, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrNotFound(err)
	}
	return &grade, nil
}

func (s *GradeService) Create(ctx domain.Context, input domain.GradeCreateInput) (*domain.Grade, error) {
	grade, err := s.gradeModel.Create(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrGradeValidationFailed(err)
	}
	return &grade, nil
}

func (s *GradeService) Update(ctx domain.Context, id string, input domain.GradeUpdateInput) (*domain.Grade, error) {
	grade, err := s.gradeModel.Update(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrGradeValidationFailed(err)
	}
	return &grade, nil
}

func (s *GradeService) Delete(ctx domain.Context, id string) error {
	_, err := s.gradeModel.Delete(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrGradeValidationFailed(err)
	}
	return nil
}

func (s *GradeService) Destroy(ctx domain.Context, id string) error {
	_, err := s.gradeModel.Destroy(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrGradeValidationFailed(err)
	}
	return nil
}
