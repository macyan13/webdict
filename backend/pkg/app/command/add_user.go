package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
)

type AddUser struct {
	Name     string
	Email    string
	Password string
}

type AddUserHandler struct {
	userRepo user.Repository
}

func (h AddUserHandler) Handle(cmd AddUser) error {
	return h.userRepo.Save(user.NewUser(cmd.Name, cmd.Email, cmd.Password))
}

func (h AddUserHandler) validate(cmd AddUser) error {
	if h.userRepo.Exist(cmd.Email) {
		return fmt.Errorf("can not create new user, a user with passed email: %s already exists", cmd.Email)
	}

	if len(cmd.Name) < 2 {
		return errors.New("can not create new user, the name must contain at least 3 character")
	}

	if len(cmd.Password) < 8 {
		return errors.New("can not create new user, the password must contain at least 3 character")
	}

	return nil
}
