package user

import "errors"

var ErrNotFound = errors.New("can not find user in store")

// Repository User domain repo
type Repository interface {
	Exist(email string) (bool, error)
	Create(user *User) error                // Create saves new user to DB
	GetByEmail(email string) (*User, error) // GetByEmail get user by email, return ErrNotFound if user not exists
}
