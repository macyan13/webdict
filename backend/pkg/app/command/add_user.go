package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
)

// AddUser create new user cmd
type AddUser struct {
	Name     string
	Email    string
	Password string
	Role     user.Role
}

// Cipher service to generate password hash before saving a user to DB
type Cipher interface {
	GenerateHash(pwd string) (string, error)
}

// AddUserHandler create new User cmd handler
type AddUserHandler struct {
	userRepo user.Repository
	cipher   Cipher
}

func NewAddUserHandler(userRepo user.Repository, cipher Cipher) AddUserHandler {
	return AddUserHandler{userRepo: userRepo, cipher: cipher}
}

// Handle performs user creation cmd
func (h AddUserHandler) Handle(cmd AddUser) (string, error) {
	hashedPwd, err := h.cipher.GenerateHash(cmd.Password)

	if err != nil {
		return "", err
	}

	if err = h.validate(cmd); err != nil {
		return "", err
	}

	u, err := user.NewUser(cmd.Name, cmd.Email, hashedPwd, cmd.Role)

	if err != nil {
		return "", err
	}

	if err := h.userRepo.Create(u); err != nil {
		return "", err
	}

	return u.ID(), nil
}

// validate checks that there is not already created user with cmd email
func (h AddUserHandler) validate(cmd AddUser) error {
	exists, err := h.userRepo.Exist(cmd.Email)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("can not create new user, a user with passed email: %s already exists", cmd.Email)
	}

	return nil
}
