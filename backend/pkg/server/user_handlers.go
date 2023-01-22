package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
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

// func (s *HTTPServer) UpdateTag() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")
// 		var request tagRequest
//
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			s.badRequest(c, fmt.Errorf("can not parse tag update request: %v", err))
// 			return
// 		}
//
// 		user, err := s.authHandler.UserFromContext(c)
// 		if err != nil {
// 			s.unauthorized(c, err)
// 		}
//
// 		if err := s.app.Commands.UpdateTag.Handle(command.UpdateTag{
// 			TagID:    c.Param(userIDParam),
// 			Tag:      request.Tag,
// 			AuthorID: user.ID,
// 		}); err != nil {
// 			s.badRequest(c, fmt.Errorf("can not Update Existing tag: %v", err))
// 			return
// 		}
//
// 		c.JSON(http.StatusOK, http.NoBody)
// 	}
// }
//
// func (s *HTTPServer) GetTagByID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")
//
// 		user, err := s.authHandler.UserFromContext(c)
// 		if err != nil {
// 			s.unauthorized(c, err)
// 		}
//
// 		view, err := s.app.Queries.SingleTag.Handle(query.SingleTag{
// 			ID:       c.Param(userIDParam),
// 			AuthorID: user.ID,
// 		})
//
// 		if err != nil {
// 			s.badRequest(c, fmt.Errorf("can get find requested tag - %v", err))
// 			return
// 		}
//
// 		c.JSON(http.StatusOK, s.tagViewToResponse(view))
// 	}
// }
//
// func (s *HTTPServer) DeleteTagByID() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Content-Type", "application/json")
//
// 		user, err := s.authHandler.UserFromContext(c)
// 		if err != nil {
// 			s.unauthorized(c, err)
// 			return
// 		}
//
// 		if err := s.app.Commands.DeleteTag.Handle(command.DeleteTag{
// 			ID:       c.Param(userIDParam),
// 			AuthorID: user.ID,
// 		}); err != nil {
// 			s.badRequest(c, fmt.Errorf("can not delete tag: %v", err))
// 			return
// 		}
//
// 		c.JSON(http.StatusOK, http.NoBody)
// 	}
// }

// func (s *HTTPServer) tagViewToResponse(tag query.TagView) tagResponse {
// 	return tagResponse{
// 		ID:  tag.ID,
// 		Tag: tag.Tag,
// 	}
// }

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
