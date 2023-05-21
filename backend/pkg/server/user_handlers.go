package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
)

const userIDParam = "userId"

func (s *HTTPServer) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request userRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse new user request: %v", err))
			return
		}

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if !usr.IsAdmin() {
			s.unauthorized(c, errors.New("authorized user don't have permissions for the action"))
			return
		}

		id, err := s.app.Commands.AddUser.Handle(command.AddUser{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
			Role:     user.Author,
		})

		if err != nil {
			if err == user.ErrEmailAlreadyExists {
				s.badRequest(c, fmt.Errorf("user with email %s already exists", request.Email))
				return
			}
			s.badRequest(c, fmt.Errorf("can not create new user: %v", err))
			return
		}

		c.JSON(http.StatusCreated, idResponse{ID: id})
	}
}

func (s *HTTPServer) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if !usr.IsAdmin() {
			s.unauthorized(c, errors.New("authorized user don't have permissions for the action"))
			return
		}

		users, err := s.app.Queries.AllUsers.Handle()

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not get users from DB - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.userViewsToResponse(users))
	}
}

func (s *HTTPServer) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if !usr.IsAdmin() {
			s.unauthorized(c, errors.New("authorized user don't have permissions for the action"))
			return
		}

		view, err := s.app.Queries.SingleUser.Handle(query.SingleUser{ID: c.Param(userIDParam)})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not find requested user - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.userViewToResponse(view))
	}
}

func (s *HTTPServer) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request userRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse user update request: %v", err))
			return
		}

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if err = s.app.Commands.UpdateUser.Handle(command.UpdateUser{
			ID:         c.Param(userIDParam),
			Name:       request.Name,
			Email:      request.Email,
			Password:   request.Password,
			Role:       user.Role(request.Role),
			IsAdminCMD: usr.IsAdmin(),
		}); err != nil {
			if err == user.ErrEmailAlreadyExists {
				s.badRequest(c, fmt.Errorf("user with email %s already exists", request.Email))
				return
			}
			s.badRequest(c, fmt.Errorf("can not Update Existing tag: %v", err))
			return
		}

		c.JSON(http.StatusOK, http.NoBody)
	}
}

func (s *HTTPServer) userViewsToResponse(users []query.UserView) []userResponse {
	responses := make([]userResponse, len(users))

	for i, usr := range users {
		responses[i] = userResponse{
			ID:    usr.ID,
			Name:  usr.Name,
			Email: usr.Email,
			Role:  usr.Role,
		}
	}

	return responses
}

func (s *HTTPServer) userViewToResponse(usr query.UserView) userResponse {
	return userResponse{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email,
		Role:  usr.Role,
	}
}
