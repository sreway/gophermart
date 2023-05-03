package domain

import (
	"github.com/google/uuid"
)

type (
	User struct {
		id           uuid.UUID
		login        string
		hashPassword string
	}
)

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Login() string {
	return u.login
}

func (u *User) HashPassword() string {
	return u.hashPassword
}

func NewUser(login, hashPassword string) *User {
	return &User{
		uuid.New(),
		login,
		hashPassword,
	}
}

func CreateUser(id uuid.UUID, login, hashPassword string) *User {
	return &User{
		id,
		login,
		hashPassword,
	}
}
