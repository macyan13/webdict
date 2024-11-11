package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/mail"
	"unicode/utf8"
)

type ListOptions struct {
	showTranscription bool
}

func NewListOptions(showTranscription bool) ListOptions {
	return ListOptions{showTranscription: showTranscription}
}

func (l *ListOptions) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"showTranscription": l.showTranscription,
	}
}

type Role int

const (
	Admin  Role = 1
	Author Role = 2
)

func (r Role) valid() bool {
	return r >= Admin && r <= Author
}

type User struct {
	id            string
	name          string
	email         string
	password      string
	role          Role
	defaultLangID string
	listOptions   ListOptions
}

func NewUser(name, email, password string, role Role) (*User, error) {
	u := User{
		id:       uuid.New().String(),
		name:     name,
		email:    email,
		password: password,
		role:     role,
		listOptions: ListOptions{
			showTranscription: false,
		},
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

func (u *User) DefaultLangID() string {
	return u.defaultLangID
}

func (u *User) ListOptions() ListOptions {
	return u.listOptions
}

func (u *User) ApplyChanges(name, email, passwd string, role Role, defaultLangID string, listOptions ListOptions) error {
	updated := *u
	updated.applyChanges(name, email, passwd, role, defaultLangID, listOptions)

	if err := updated.validate(); err != nil {
		return err
	}

	u.applyChanges(name, email, passwd, role, defaultLangID, listOptions)
	return nil
}

func (u *User) applyChanges(name, email, passwd string, role Role, defaultLangID string, listOptions ListOptions) {
	u.name = name
	u.email = email
	u.role = role
	u.defaultLangID = defaultLangID
	u.password = passwd
	u.listOptions = listOptions
}

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            u.id,
		"name":          u.name,
		"email":         u.email,
		"password":      u.password,
		"role":          int(u.role),
		"defaultLangID": u.defaultLangID,
		"listOptions":   u.listOptions.ToMap(),
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

	if !u.role.valid() {
		err = errors.Join(fmt.Errorf("invalid user role passed - %d", u.role), err)
	}

	return err
}

func UnmarshalFromDB(
	id string,
	name string,
	email string,
	password string,
	role Role,
	defaultLangID string,
	listOptions ListOptions,
) *User {
	return &User{
		id:            id,
		name:          name,
		email:         email,
		password:      password,
		role:          role,
		defaultLangID: defaultLangID,
		listOptions:   listOptions,
	}
}
