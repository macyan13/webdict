package user

import (
	"errors"
	"github.com/google/uuid"
)

type User struct {
	id       string
	name     string
	email    string
	password string
	role     Role
}

type Role int

const (
	Admin  Role = 1
	Author Role = 2
)

func NewUser(name, email, password string, role Role) (User, error) {
	u := User{
		id:       uuid.New().String(),
		name:     name,
		email:    email,
		password: password,
		role:     role,
	}

	if err := u.validate(); err != nil {
		return User{}, err
	}

	return u, nil
}

func (u *User) ID() string {
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

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       u.id,
		"name":     u.name,
		"email":    u.email,
		"password": u.password,
		"role":     int(u.role),
	}
}

func UnmarshalFromDB(
	id string,
	name string,
	email string,
	password string,
	role int,
) User {
	return User{
		id:       id,
		name:     name,
		email:    email,
		password: password,
		role:     Role(role),
	}
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
