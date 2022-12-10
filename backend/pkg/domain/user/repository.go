package user

import "errors"

var NotFoundErr = errors.New("can not find user in storage")

type Repository interface {
	Exist(email string) (bool, error)
	Create(user *User) error
	GetByEmail(email string) (*User, error)
}
