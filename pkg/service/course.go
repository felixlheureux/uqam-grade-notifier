package service

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/model"
)

type CourseService struct {
	courseModel *model.Course
}

func NewCourseService(courseModel *model.Course) *CourseService {
	return &CourseService{
		courseModel: courseModel,
	}
}

func (s *CourseService) Get(ctx domain.Context) ([]domain.Course, error) {
	courses, err := s.courseModel.Get("", ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrCourseStorageFailed(err)
	}
	return []domain.Course{courses}, nil
}

func (s *CourseService) FindOne(ctx domain.Context, id string) (*domain.Course, error) {
	course, err := s.courseModel.Get(id, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrNotFound(err)
	}
	return &course, nil
}

func (s *CourseService) Create(ctx domain.Context, input domain.CourseCreateInput) (*domain.Course, error) {
	course, err := s.courseModel.Create(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrCourseStorageFailed(err)
	}
	return &course, nil
}

func (s *CourseService) Update(ctx domain.Context, id string, input domain.CourseUpdateInput) (*domain.Course, error) {
	course, err := s.courseModel.Update(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrCourseStorageFailed(err)
	}
	return &course, nil
}

func (s *CourseService) Delete(ctx domain.Context, id string) error {
	_, err := s.courseModel.Delete(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrCourseStorageFailed(err)
	}
	return nil
}

func (s *CourseService) Destroy(ctx domain.Context, id string) error {
	_, err := s.courseModel.Destroy(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrCourseStorageFailed(err)
	}
	return nil
}
