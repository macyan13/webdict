package mongo

import (
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo_fromDomainToModel(t *testing.T) {
	name := "John"
	email := "John@do.com"
	password := "12345678"
	role := user.Admin
	usr, err := user.NewUser(name, email, password, role)

	assert.Nil(t, err)
	repo := UserRepo{}

	model, err := repo.fromDomainToModel(*usr)
	assert.Nil(t, err)
	assert.Equal(t, usr.Id(), model.Id)
	assert.Equal(t, name, name)
	assert.Equal(t, email, model.Email)
	assert.Equal(t, password, model.Password)
	assert.Equal(t, int(role), model.Role)
}
