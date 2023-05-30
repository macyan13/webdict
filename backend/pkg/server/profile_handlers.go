package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
)

func (s *HTTPServer) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		view, err := s.app.Queries.SingleUser.Handle(query.SingleUser{ID: usr.ID})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not find requested user - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.userViewToResponse(view))
	}
}

func (s *HTTPServer) UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request updateProfileRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse profile update request: %v", err))
			return
		}

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if err = s.app.Commands.UpdateProfile.Handle(command.UpdateProfile{
			ID:              usr.ID,
			Name:            request.Name,
			Email:           request.Email,
			CurrentPassword: request.CurrentPassword,
			NewPassword:     request.NewPassword,
			DefaultLangID:   request.DefaultLangID,
		}); err != nil {
			if err == user.ErrEmailAlreadyExists {
				s.badRequest(c, fmt.Errorf("user with email %s already exists", request.Email))
				return
			}
			s.badRequest(c, fmt.Errorf("can not update user: %v", err))
			return
		}

		c.JSON(http.StatusOK, http.NoBody)
	}
}
