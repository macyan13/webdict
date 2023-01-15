package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/mail"
	"unicode/utf8"
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

func NewUser(name, email, password string, role Role) (*User, error) {
	u := User{
		id:       uuid.New().String(),
		name:     name,
		email:    email,
		password: password,
		role:     role,
	}

	if err := u.validate(); err != nil {
		return nil, err
	}

	return &u, nil
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
) *User {
	return &User{
		id:       id,
		name:     name,
		email:    email,
		password: password,
		role:     Role(role),
	}
}

func (u *User) validate() error {
	nameCount := utf8.RuneCountInString(u.name)

	if nameCount < 2 {
		return fmt.Errorf("name must contain at least 2 characters, %d passed (%s)", nameCount, u.name)
	}

	if nameCount > 30 {
		return fmt.Errorf("name max size is 30 characters, %d passed (%s)", nameCount, u.name)
	}

	if _, err := mail.ParseAddress(u.email); err != nil {
		return fmt.Errorf("email is not valid: %s", err.Error())
	}

	// it should never happen as domain receives passwd hash from cipher
	if len(u.password) < 8 {
		return errors.New("password must contain at least 8 character")
	}

	return nil
}
