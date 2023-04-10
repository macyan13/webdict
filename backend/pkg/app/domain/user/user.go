package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
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
	var result error
	nameCount := utf8.RuneCountInString(u.name)

	if nameCount < 2 {
		result = multierror.Append(result, fmt.Errorf("name must contain at least 2 characters, %d passed (%s)", nameCount, u.name))
	}

	if nameCount > 30 {
		result = multierror.Append(result, fmt.Errorf("name max size is 30 characters, %d passed (%s)", nameCount, u.name))
	}

	if _, err := mail.ParseAddress(u.email); err != nil {
		result = multierror.Append(result, fmt.Errorf("email is not valid: %s", err.Error()))
	}

	// it should never happen as domain receives passwd as hash from cipher
	if len(u.password) < 8 {
		result = multierror.Append(result, errors.New("password must contain at least 8 character"))
	}

	return result
}