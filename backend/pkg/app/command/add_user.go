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

type Cipher interface {
	GenerateHash(pwd string) (string, error)
}

type AddUserHandler struct {
	userRepo user.Repository
	cipher   Cipher
}

func NewAddUserHandler(userRepo user.Repository, cipher Cipher) AddUserHandler {
	return AddUserHandler{userRepo: userRepo, cipher: cipher}
}

func (h AddUserHandler) Handle(cmd AddUser) error {
	hashedPwd, err := h.cipher.GenerateHash(cmd.Password)

	if err != nil {
		return err
	}

	u, err := user.NewUser(cmd.Name, cmd.Email, hashedPwd, user.Admin)

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
