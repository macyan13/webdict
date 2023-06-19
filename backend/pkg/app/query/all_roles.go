package query

import "github.com/macyan13/webdict/backend/pkg/app/domain/user"

type AllRolesHandler struct {
	mapper *RoleConverter
	roles  []user.Role
}

func NewAllRolesHandler() AllRolesHandler {
	return AllRolesHandler{mapper: NewRoleMapper(), roles: []user.Role{user.Admin, user.Author}}
}

func (h AllRolesHandler) Handle() ([]RoleView, error) {
	mappedRoles := make([]RoleView, 0, len(h.roles))

	for _, role := range h.roles {
		roleView, err := h.mapper.RoleToView(role)
		if err != nil {
			return []RoleView{}, err
		}

		mappedRoles = append(mappedRoles, roleView)
	}

	return mappedRoles, nil
}
