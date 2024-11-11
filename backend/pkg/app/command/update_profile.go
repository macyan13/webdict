package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
)

type UpdateProfile struct {
	ID              string
	Name            string
	Email           string
	CurrentPassword string
	NewPassword     string
	DefaultLangID   string
	ListOptions     user.ListOptions
}

type UpdateProfileHandler struct {
	userRepo user.Repository
	cipher   Cipher
	langRepo lang.Repository
}

func NewUpdateProfileHandler(userRepo user.Repository, cipher Cipher, langRepo lang.Repository) UpdateProfileHandler {
	return UpdateProfileHandler{userRepo: userRepo, cipher: cipher, langRepo: langRepo}
}

func (h UpdateProfileHandler) Handle(cmd UpdateProfile) error {
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

	if err = usr.ApplyChanges(cmd.Name, cmd.Email, passwd, usr.Role(), cmd.DefaultLangID, cmd.ListOptions); err != nil {
		return err
	}

	return h.userRepo.Update(usr)
}

func (h UpdateProfileHandler) processPasswd(cmd UpdateProfile, userHash string) (string, error) {
	if cmd.NewPassword == "" {
		return userHash, nil
	}

	if !h.cipher.ComparePasswords(userHash, cmd.CurrentPassword) {
		return "", errors.New("current password is not valid")
	}

	return h.cipher.GenerateHash(cmd.NewPassword)
}

func (h UpdateProfileHandler) validate(cmd UpdateProfile) error {
	if cmd.DefaultLangID == "" {
		return nil
	}

	exist, err := h.langRepo.Exist(cmd.DefaultLangID, cmd.ID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("lang with id: %s is not found", cmd.DefaultLangID)
	}

	return nil
}
