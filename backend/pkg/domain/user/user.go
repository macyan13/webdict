package user

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	id        string
	name      string
	email     string
	password  string
	createdAt time.Time
	role      Role
}

type Role int

const (
	Admin  Role = 1
	Author Role = 2
)

func NewUser(name, email, password string, role Role) (*User, error) {
	u := User{
		id:        uuid.New().String(),
		name:      name,
		email:     email,
		password:  password,
		createdAt: time.Now(),
		role:      role,
	}

	if err := u.validate(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) validate() error {
	// todo: add validation for email
	if len(u.name) < 3 {
		return errors.New("can not create new user, the name must contain at least 3 character")
	}

	if len(u.password) < 8 {
		return errors.New("can not create new user, the password must contain at least 3 character")
	}

	return nil
}
