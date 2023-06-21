package command

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
)

type UpdateUser struct {
	ID       string
	Name     string
	Email    string
	Password string
	Role     user.Role
}

type UpdateUserHandler struct {
	userRepo user.Repository
	cipher   Cipher
}

func NewUpdateUserHandler(userRepo user.Repository, cipher Cipher) UpdateUserHandler {
	return UpdateUserHandler{userRepo: userRepo, cipher: cipher}
}

func (h UpdateUserHandler) Handle(cmd UpdateUser) error {
	usr, err := h.userRepo.Get(cmd.ID)
	if err != nil {
		return err
	}

	passwd, err := h.processPasswd(cmd, usr.Password())
	if err != nil {
		return err
	}

	if err = usr.ApplyChanges(cmd.Name, cmd.Email, passwd, cmd.Role, usr.DefaultLangID()); err != nil {
		return err
	}

	return h.userRepo.Update(usr)
}

func (h UpdateUserHandler) processPasswd(cmd UpdateUser, userHash string) (string, error) {
	if cmd.Password == "" {
		return userHash, nil
	}

	return h.cipher.GenerateHash(cmd.Password)
}
