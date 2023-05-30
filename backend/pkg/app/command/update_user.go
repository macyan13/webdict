package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
)

// UpdateUser update existing user cmd
type UpdateUser struct {
	ID              string
	Name            string
	Email           string
	CurrentPassword string
	NewPassword     string
	Role            user.Role
	IsAdminCMD      bool
	DefaultLangID   string
}

// UpdateUserHandler update existing user cmd handler
type UpdateUserHandler struct {
	userRepo user.Repository
	cipher   Cipher
	langRepo lang.Repository
}

func NewUpdateUserHandler(userRepo user.Repository, cipher Cipher, langRepo lang.Repository) UpdateUserHandler {
	return UpdateUserHandler{userRepo: userRepo, cipher: cipher, langRepo: langRepo}
}

// Handle applies cmd changes to tag and saves it to DB
func (h UpdateUserHandler) Handle(cmd UpdateUser) error {
	usr, err := h.userRepo.Get(cmd.ID)
	if err != nil {
		return err
	}

	if err = h.validate(cmd); err != nil {
		return err
	}

	passwd, err := h.processPasswd(cmd, usr.Password())
	if err != nil {
		return err
	}

	role, err := h.processRole(cmd, usr)
	if err != nil {
		return err
	}

	if err = usr.ApplyChanges(cmd.Name, cmd.Email, passwd, role, cmd.DefaultLangID); err != nil {
		return err
	}

	return h.userRepo.Update(usr)
}

func (h UpdateUserHandler) processPasswd(cmd UpdateUser, userHash string) (string, error) {
	if cmd.NewPassword == "" {
		return userHash, nil
	}

	if !h.cipher.ComparePasswords(userHash, cmd.CurrentPassword) {
		return "", errors.New("current password is not valid")
	}

	return h.cipher.GenerateHash(cmd.NewPassword)
}

func (h UpdateUserHandler) processRole(cmd UpdateUser, usr *user.User) (user.Role, error) {
	if cmd.IsAdminCMD {
		return cmd.Role, nil
	}

	if cmd.Role != user.NotSet && usr.Role() != cmd.Role {
		return 0, errors.New("attempt to change the role from not admin cmd")
	}

	return usr.Role(), nil
}

func (h UpdateUserHandler) validate(cmd UpdateUser) error {
	if cmd.DefaultLangID == "" {
		return nil
	}

	_, err := h.langRepo.Get(cmd.DefaultLangID, cmd.ID)
	return err
}
