package command

import (
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

func NewAddUserHandler(userRepo user.Repository) AddUserHandler {
	return AddUserHandler{userRepo: userRepo}
}

func (h AddUserHandler) Handle(cmd AddUser) error {
	u, err := user.NewUser(cmd.Name, cmd.Email, cmd.Password)

	if err != nil {
		return err
	}

	return h.userRepo.Save(u)
}

func (h AddUserHandler) validate(cmd AddUser) error {
	if h.userRepo.Exist(cmd.Email) {
		return fmt.Errorf("can not create new user, a user with passed email: %s already exists", cmd.Email)
	}

	return nil
}
