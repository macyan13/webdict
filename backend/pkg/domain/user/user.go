package user

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

func newUser(name, email, password string) *User {
	return &User{
		Id:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
