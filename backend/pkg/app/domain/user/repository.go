package user

import "errors"

var ErrNotFound = errors.New("can not find user in store")

// Repository User domain repo
type Repository interface {
	Exist(email string) (bool, error)       // Exist checks if player with the email exists
	Create(user *User) error                // Create saves new user to DB
	GetByEmail(email string) (*User, error) // GetByEmail gets user by email, return ErrNotFound if user not exists
	Get(id string) (*User, error)           // Get gets user by id, return ErrNotFound if user not exists
	Update(usr *User) error                 // Update saves the updated usr entity to store
}
