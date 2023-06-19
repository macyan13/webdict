package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
)

func (s *HTTPServer) GetRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		roles, err := s.app.Queries.AllRoles.Handle()
		if err != nil {
			s.badRequest(c, fmt.Errorf("can not provide role list: %s", err.Error()))
			return
		}

		c.JSON(http.StatusOK, s.roleViewsToResponse(roles))
	}
}

func (s *HTTPServer) roleViewsToResponse(roles []query.RoleView) rolesResponse {
	roleViews := make([]roleResponse, len(roles))

	for i, role := range roles {
		roleViews[i] = s.roleViewToResponse(role)
	}

	return rolesResponse{Roles: roleViews}
}

func (s *HTTPServer) roleViewToResponse(role query.RoleView) roleResponse {
	return roleResponse{
		ID:      role.ID,
		Name:    role.Name,
		IsAdmin: role.IsAdmin,
	}
}
