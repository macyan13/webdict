package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
)

const langIDParam = "langId"

func (s *HTTPServer) CreateLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request langRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse new lang request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		id, err := s.app.Commands.AddLang.Handle(command.AddLang{
			Name:     request.Name,
			AuthorID: user.ID,
		})

		if err == command.ErrLangAlreadyExists {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("lang %s already exists", request.Name))
			return
		}

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not create new lang: %v", err))
			return
		}

		c.JSON(http.StatusCreated, idResponse{ID: id})
	}
}

func (s *HTTPServer) UpdateLang() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request langRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse lang update request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if err := s.app.Commands.UpdateLang.Handle(command.UpdateLang{
			ID:       c.Param(langIDParam),
			Name:     request.Name,
			AuthorID: user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not Update Existing lang: %v", err))
			return
		}

		c.JSON(http.StatusOK, http.NoBody)
	}
}

func (s *HTTPServer) DeleteLangByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if err := s.app.Commands.DeleteLang.Handle(command.DeleteLang{
			ID:       c.Param(langIDParam),
			AuthorID: user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not delete lang: %v", err))
			return
		}

		c.JSON(http.StatusOK, http.NoBody)
	}
}

func (s *HTTPServer) GetLangByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		view, err := s.app.Queries.SingleLang.Handle(query.SingleLang{
			ID:       c.Param(langIDParam),
			AuthorID: user.ID,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not find requested lang - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.langViewToResponse(view))
	}
}

func (s *HTTPServer) GetLangs() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		usr, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		views, err := s.app.Queries.AllLangs.Handle(query.AllLangs{AuthorID: usr.ID})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not get langs from DB - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.langViewsToResponse(views))
	}
}

func (s *HTTPServer) langViewToResponse(ln query.LangView) langResponse {
	return langResponse{
		ID:   ln.ID,
		Name: ln.Name,
	}
}

func (s *HTTPServer) langViewsToResponse(langs []query.LangView) []langResponse {
	responses := make([]langResponse, len(langs))

	for i, l := range langs {
		responses[i] = s.langViewToResponse(l)
	}

	return responses
}
