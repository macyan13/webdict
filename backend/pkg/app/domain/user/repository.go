package user

import "errors"

var ErrNotFound = errors.New("can not find user in store")
var ErrEmailAlreadyExists = errors.New("user with such email already exists")

// Repository User domain repo
type Repository interface {
	Create(user *User) error                // Create saves new user to DB, return ErrEmailAlreadyExists when user with email already exists
	GetByEmail(email string) (*User, error) // GetByEmail gets user by email, return ErrNotFound if user not exists
	Get(id string) (*User, error)           // Get gets user by id, return ErrNotFound if user not exists
	Update(usr *User) error                 // Update saves the updated usr entity to store, return ErrEmailAlreadyExists when user with email already exists
	Delete(id string) (int, error)          // Delete removes user from DB
}
