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
}

func NewUser(name, email, password string) (*User, error) {
	u := User{
		id:        uuid.New().String(),
		name:      name,
		email:     email,
		password:  password,
		createdAt: time.Now(),
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
