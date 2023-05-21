package server //nolint:dupl // it's not fully duplicate

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
)

const tagIDParam = "tagId"

func (s *HTTPServer) CreateTag() gin.HandlerFunc { //nolint:dupl // it's not fully duplicate
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
			return
		}

		id, err := s.app.Commands.AddTag.Handle(command.AddTag{
			Tag:      request.Tag,
			AuthorID: user.ID,
		})

		if err == tag.ErrTagAlreadyExists {
			s.badRequest(c, fmt.Errorf("tag %s already exists", request.Tag))
			return
		}

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not create new tag: %v", err))
			return
		}

		c.JSON(http.StatusCreated, idResponse{ID: id})
	}
}

func (s *HTTPServer) GetTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		tags, err := s.app.Queries.AllTags.Handle(query.AllTags{AuthorID: user.ID})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not get tags from DB - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.tagViewsToResponse(tags))
	}
}

func (s *HTTPServer) UpdateTag() gin.HandlerFunc { //nolint:dupl // it's not fully duplicate
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
			return
		}

		if err = s.app.Commands.UpdateTag.Handle(command.UpdateTag{
			TagID:    c.Param(tagIDParam),
			Tag:      request.Tag,
			AuthorID: user.ID,
		}); err != nil {
			if err == tag.ErrTagAlreadyExists {
				s.badRequest(c, fmt.Errorf("tag %s already exists", request.Tag))
				return
			}
			s.badRequest(c, fmt.Errorf("can not Update Existing tag: %v", err))
			return
		}

		c.JSON(http.StatusOK, http.NoBody)
	}
}

func (s *HTTPServer) GetTagByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		view, err := s.app.Queries.SingleTag.Handle(query.SingleTag{
			ID:       c.Param(tagIDParam),
			AuthorID: user.ID,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not find requested tag - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.tagViewToResponse(view))
	}
}

func (s *HTTPServer) DeleteTagByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		if err := s.app.Commands.DeleteTag.Handle(command.DeleteTag{
			ID:       c.Param(tagIDParam),
			AuthorID: user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not delete tag: %v", err))
			return
		}

		c.JSON(http.StatusOK, http.NoBody)
	}
}

func (s *HTTPServer) tagViewToResponse(tg query.TagView) tagResponse {
	return tagResponse{
		ID:  tg.ID,
		Tag: tg.Tag,
	}
}

func (s *HTTPServer) tagViewsToResponse(tags []query.TagView) []tagResponse {
	responses := make([]tagResponse, len(tags))

	for i, t := range tags {
		responses[i] = s.tagViewToResponse(t)
	}

	return responses
}
