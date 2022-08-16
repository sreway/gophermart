package user

import (
	"context"
	"errors"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/utils"
)

type User struct {
	repo usecase.UserRepo
}

func New(repo usecase.UserRepo) *User {
	return &User{
		repo: repo,
	}
}

func (uc *User) Register(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return entity.NewUserError(user.Login, entity.ErrUserIncorrectData)
	}

	byteHash, err := utils.HashAndSalt(user.Password)
	if err != nil {
		return entity.ErrUserIncorrectData
	}

	user.PasswordHash = string(byteHash)

	err = uc.repo.Add(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *User) Login(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := user.Validate(); err != nil {
		return nil, entity.ErrUserIncorrectData
	}

	existUser, err := uc.repo.Get(ctx, user.Login)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return nil, entity.ErrUserNotFound
		}
		return user, err
	}

	if !utils.ComparePassword(existUser.PasswordHash, user.Password) {
		return nil, entity.ErrUserIncorrectPassword
	}
	return existUser, nil
}
