package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
)

// AddUser create new user cmd
type AddUser struct {
	Name     string
	Email    string
	Password string
	Role     user.Role
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
	if !cmd.Role.Valid() {
		return "", fmt.Errorf("attempt to create user with invalid role: %d", cmd.Role)
	}

	hashedPwd, err := h.cipher.GenerateHash(cmd.Password)

	if err != nil {
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
