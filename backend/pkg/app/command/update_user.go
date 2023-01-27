package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
)

// UpdateUser update existing user cmd
type UpdateUser struct {
	ID         string
	Name       string
	Email      string
	Password   string
	Role       user.Role
	IsAdminCMD bool
}

// UpdateUserHandler update existing user cmd handler
type UpdateUserHandler struct {
	userRepo user.Repository
	cipher   Cipher
}

func NewUpdateUserHandler(userRepo user.Repository, cipher Cipher) UpdateUserHandler {
	return UpdateUserHandler{userRepo: userRepo, cipher: cipher}
}

// Handle applies cmd changes to tag and saves it to DB
func (h UpdateUserHandler) Handle(cmd UpdateUser) error {
	if !cmd.Role.Valid() {
		return fmt.Errorf("can not update user, invalid role passed - %d", cmd.Role)
	}

	usr, err := h.userRepo.Get(cmd.ID)
	if err != nil {
		return err
	}

	if usr.Email() != cmd.Email {
		free, err := h.emailFree(cmd)

		if err != nil {
			return err
		}

		if !free {
			return fmt.Errorf("email %s already in use", cmd.Email)
		}
	}

	if usr.Role() != cmd.Role && !cmd.IsAdminCMD {
		return errors.New("can not update user, attempt to change the role from not admin cmd")
	}

	var passwd string
	if cmd.Password != "" {
		passwd, err = h.cipher.GenerateHash(cmd.Password)
		if err != nil {
			return err
		}
	}

	if err := usr.ApplyChanges(cmd.Name, cmd.Email, passwd, cmd.Role); err != nil {
		return err
	}

	return h.userRepo.Update(usr)
}

func (h UpdateUserHandler) emailFree(cmd UpdateUser) (bool, error) {
	_, err := h.userRepo.GetByEmail(cmd.Email)

	if err == user.ErrNotFound {
		return true, nil
	}

	return false, err
}
