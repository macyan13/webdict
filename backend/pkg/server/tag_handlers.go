package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
)

const tagIDParam = "tagId"

func (s *HTTPServer) CreateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request tagRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse new tag request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if err := s.app.Commands.AddTag.Handle(command.AddTag{
			Tag:      request.Tag,
			AuthorID: user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not create new tag: %v", err))
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *HTTPServer) GetTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		tag, err := s.app.Queries.AllTags.Handle(query.AllTags{AuthorID: user.ID})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not get tags from DB - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.tagViewsToResponse(tag))
	}
}

func (s *HTTPServer) UpdateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request tagRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse tag update request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if err := s.app.Commands.UpdateTag.Handle(command.UpdateTag{
			TagID:    c.Param(tagIDParam),
			Tag:      request.Tag,
			AuthorID: user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not Update Existing tag: %v", err))
			return
		}

		response := map[string]any{
			"status": "success",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *HTTPServer) GetTagByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		view, err := s.app.Queries.SingleTag.Handle(query.SingleTag{
			ID:       c.Param(tagIDParam),
			AuthorID: user.ID,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can get find requested tag - %v", err))
		}

		c.JSON(http.StatusOK, s.tagModelToResponse(view))
	}
}

func (s *HTTPServer) DeleteTagByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if err := s.app.Commands.DeleteTag.Handle(command.DeleteTag{
			ID:       c.Param(tagIDParam),
			AuthorID: user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not delete tag: %v", err))
			return
		}

		response := map[string]any{
			"status": "success",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *HTTPServer) tagModelToResponse(tag query.TagView) tagResponse {
	return tagResponse{
		ID:  tag.ID,
		Tag: tag.Tag,
	}
}

func (s *HTTPServer) tagViewsToResponse(tags []query.TagView) []tagResponse {
	responses := make([]tagResponse, len(tags))

	for i, t := range tags {
		responses[i] = tagResponse{
			ID:  t.ID,
			Tag: t.Tag,
		}
	}

	return responses
}
