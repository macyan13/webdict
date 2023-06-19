package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
)

type UserRepo struct {
	storage       map[string]*user.User
	roleConverter *query.RoleConverter
}

func NewUserRepository(roleMapper *query.RoleConverter) *UserRepo {
	return &UserRepo{
		storage:       map[string]*user.User{},
		roleConverter: roleMapper,
	}
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

func (u *UserRepo) Get(id string) (*user.User, error) {
	usr, ok := u.storage[id]

	if ok {
		return usr, nil
	}

	return nil, fmt.Errorf("not found")
}

func (u *UserRepo) Update(usr *user.User) error {
	u.storage[usr.ID()] = usr
	return nil
}

func (u *UserRepo) GetAllViews() ([]query.UserView, error) {
	results := make([]query.UserView, 0, len(u.storage))

	for s := range u.storage {
		userData := u.storage[s].ToMap()
		role, err := u.roleConverter.RoleToView(user.Role(userData["role"].(int)))
		if err != nil {
			return []query.UserView{}, err
		}

		results = append(results, query.UserView{
			ID:    userData["id"].(string),
			Name:  userData["name"].(string),
			Email: userData["email"].(string),
			Role:  role,
		})
	}

	return results, nil
}

func (u *UserRepo) GetView(id string) (query.UserView, error) {
	t, ok := u.storage[id]

	if ok {
		userData := t.ToMap()
		role, err := u.roleConverter.RoleToView(user.Role(userData["role"].(int)))
		if err != nil {
			return query.UserView{}, err
		}

		return query.UserView{
			ID:    userData["id"].(string),
			Name:  userData["name"].(string),
			Email: userData["email"].(string),
			Role:  role,
		}, nil
	}

	return query.UserView{}, fmt.Errorf("not found")
}
