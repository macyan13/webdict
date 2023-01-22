package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
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

func (u *UserRepo) Create(usr *user.User) error {
	u.storage[usr.ID()] = usr
	return nil
}

func (u *UserRepo) GetByEmail(email string) (*user.User, error) {
	for _, el := range u.storage {
		if el.Email() == email {
			return el, nil
		}
	}

	return nil, user.ErrNotFound
}

func (u *UserRepo) GetAllViews() ([]query.UserView, error) {
	results := make([]query.UserView, 0, len(u.storage))

	for s := range u.storage {
		userData := u.storage[s].ToMap()
		results = append(results, query.UserView{
			ID:    userData["id"].(string),
			Name:  userData["name"].(string),
			Email: userData["email"].(string),
			Role:  userData["role"].(int),
		})
	}

	return results, nil
}

func (u *UserRepo) GetView(id string) (query.UserView, error) {
	t, ok := u.storage[id]

	if ok {
		userData := t.ToMap()
		return query.UserView{
			ID:    userData["id"].(string),
			Name:  userData["name"].(string),
			Email: userData["email"].(string),
			Role:  userData["role"].(int),
		}, nil
	}

	return query.UserView{}, fmt.Errorf("not found")
}
