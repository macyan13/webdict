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

func (r Role) Valid() bool {
	return r >= Admin && r <= Author
}

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

func (u *User) ApplyChanges(name, email, passwd string, role Role) error {
	updated := *u
	updated.applyChanges(name, email, passwd, role)

	if err := updated.validate(); err != nil {
		return err
	}

	u.applyChanges(name, email, passwd, role)
	return nil
}

func (u *User) applyChanges(name, email, passwd string, role Role) {
	u.name = name
	u.email = email
	u.role = role

	if passwd != "" {
		u.password = passwd
	}
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
	var err error
	nameCount := utf8.RuneCountInString(u.name)

	if nameCount < 2 {
		err = errors.Join(fmt.Errorf("name must contain at least 2 characters, %d passed (%s)", nameCount, u.name), err)
	}

	if nameCount > 30 {
		err = errors.Join(fmt.Errorf("name max size is 30 characters, %d passed (%s)", nameCount, u.name), err)
	}

	if _, addressErr := mail.ParseAddress(u.email); addressErr != nil {
		err = errors.Join(fmt.Errorf("email is not valid: %s", addressErr.Error()), err)
	}

	// it should never happen as domain receives passwd as hash from cipher
	if len(u.password) < 8 {
		err = errors.Join(errors.New("password must contain at least 8 character"), err)
	}

	return err
}
