package model

import (
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/uptrace/bun"
)

type User struct {
	BaseModel

	ID           string `bun:",pk"`
	Email        string
	UQAMPassword string
	UpdatedAt    time.Time
	CreatedAt    time.Time

	model[User, domain.User, domain.UserCreateInput, domain.UserUpdateInput, domain.UserFilters]
}

func NewUser(db *bun.DB) *User {
	user := &User{}
	user.SetDB(db)
	return user
}

func (u User) id() string {
	return u.ID
}

func (User) idPrefix() string {
	return "usr_"
}

func (User) tableName() string {
	return "users"
}

func (User) postprocess(model User) domain.User {
	return domain.User{
		ID:           model.ID,
		Email:        model.Email,
		UQAMPassword: model.UQAMPassword,
		UpdatedAt:    model.UpdatedAt,
		CreatedAt:    model.CreatedAt,
	}
}
