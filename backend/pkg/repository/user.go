package repository

import "github.com/macyan13/webdict/backend/pkg/domain/user"

type userRepo struct {
	storage map[string]user.User
}

func NewUserRepository() *userRepo {
	return &userRepo{
		storage: map[string]user.User{},
	}
}

func (u *userRepo) Exist(id string) bool {
	_, ok := u.storage[id]
	return ok
}

func (u userRepo) Save(user *user.User) error {
	u.storage[user.Id] = *user
	return nil
}

func (u userRepo) GetByEmail(email string) *user.User {
	for _, el := range u.storage {
		if el.Email == email {
			return &el
		}
	}

	return nil
}
