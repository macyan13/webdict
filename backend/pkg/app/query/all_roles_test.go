package query

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllRolesHandler_Handle_PositiveCase(t *testing.T) {
	handler := NewAllRolesHandler()
	roles, err := handler.Handle()
	assert.Nil(t, err)
	assert.Equal(t, len(handler.roles), len(roles))
}

func TestAllRolesHandler_Handle_UnsupportedRole(t *testing.T) {
	handler := NewAllRolesHandler()
	handler.roles = []user.Role{user.Role(3)}
	_, err := handler.Handle()
	assert.Error(t, err)
}
