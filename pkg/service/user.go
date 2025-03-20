package service

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/model"
)

type UserService struct {
	userModel *model.User
}

func NewUserService(userModel *model.User) *UserService {
	return &UserService{
		userModel: userModel,
	}
}

func (s *UserService) Get(ctx domain.Context) ([]domain.User, error) {
	users, err := s.userModel.Get("", ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUserGetFailed(err)
	}
	return []domain.User{users}, nil
}

func (s *UserService) FindOne(ctx domain.Context, id string) (*domain.User, error) {
	user, err := s.userModel.Get(id, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUserFindOneFailed(err)
	}
	return &user, nil
}

func (s *UserService) Create(ctx domain.Context, input domain.UserCreateInput) (*domain.User, error) {
	user, err := s.userModel.Create(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUserCreateFailed(err)
	}
	return &user, nil
}

func (s *UserService) Update(ctx domain.Context, id string, input domain.UserUpdateInput) (*domain.User, error) {
	user, err := s.userModel.Update(input, ctx.Request().Context(), nil)
	if err != nil {
		return nil, domain.ErrUserUpdateFailed(err)
	}
	return &user, nil
}

func (s *UserService) Delete(ctx domain.Context, id string) error {
	_, err := s.userModel.Delete(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrUserDeleteFailed(err)
	}
	return nil
}

func (s *UserService) Destroy(ctx domain.Context, id string) error {
	_, err := s.userModel.Destroy(id, ctx.Request().Context(), nil)
	if err != nil {
		return domain.ErrUserDestroyFailed(err)
	}
	return nil
}
