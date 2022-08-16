package entity

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID           uint   `json:"-"`
	Login        string `json:"login" db:"login" validate:"required"`
	Password     string `json:"password" validate:"required"`
	PasswordHash string `json:"-" db:"password_hash"`
}

func (user *User) Validate() error {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return err
	}
	return nil
}
