package user

import "errors"

var NotFoundErr = errors.New("can not find user in storage")

// Repository User domain repo
type Repository interface {
	Exist(email string) (bool, error)
	Create(user User) error                // Create saves new user to DB
	GetByEmail(email string) (User, error) // GetByEmail get user by email, return NotFoundErr if user not exists
}
