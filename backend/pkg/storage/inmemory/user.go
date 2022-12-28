package inmemory

import (
	"github.com/macyan13/webdict/backend/pkg/domain/user"
)

type UserRepo struct {
	storage map[string]*user.User
}

func NewUserRepository() *UserRepo {
	return &UserRepo{
		storage: map[string]*user.User{},
	}
}

func (u *UserRepo) Exist(email string) (bool, error) {
	_, err := u.GetByEmail(email)

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (u *UserRepo) Create(user user.User) error {
	u.storage[user.Id()] = &user
	return nil
}

func (u *UserRepo) GetByEmail(email string) (user.User, error) {
	for _, el := range u.storage {
		if el.Email() == email {
			return *el, nil
		}
	}

	return user.User{}, user.NotFoundErr
}
