package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"log"
	"net/http"
)

const tagIdParam = "tagId"

func (s *HttpServer) CreateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request tagRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse new Tag request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id

		if err := s.app.Commands.AddTag.Handle(command.AddTag{
			Tag:      request.Tag,
			AuthorId: authorId,
		}); err != nil {
			log.Printf("[Error] Can not create new Tag: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *HttpServer) GetTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id
		c.JSON(http.StatusOK, s.tagModelsToResponse(s.app.Queries.AllTags.Handle(query.AllTags{AuthorId: authorId})))
	}
}

func (s *HttpServer) UpdateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request tagRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse Tag update request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id

		if err := s.app.Commands.UpdateTag.Handle(command.UpdateTag{
			TagId:    c.Param(tagIdParam),
			Tag:      request.Tag,
			AuthorId: authorId,
		}); err != nil {
			log.Printf("[Error] Can not Update Existing Tag: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		response := map[string]any{
			"status": "success",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *HttpServer) GetTagById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id

		record := s.app.Queries.SingleTag.Handle(query.SingleTag{
			Id:       c.Param(tagIdParam),
			AuthorId: authorId,
		})

		if record == nil {
			c.JSON(http.StatusBadRequest, "Can not find requested tag")
		}

		c.JSON(http.StatusOK, s.tagModelToResponse(*record))
	}
}

func (s *HttpServer) DeleteTagById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id

		if err := s.app.Commands.DeleteTag.Handle(command.DeleteTag{
			Id:       c.Param(tagIdParam),
			AuthorId: authorId,
		}); err != nil {
			log.Printf("[Error] Can not delete tag: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		response := map[string]any{
			"status": "success",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *HttpServer) tagModelToResponse(tag query.Tag) tagResponse {
	return tagResponse{
		Id:  tag.Id,
		Tag: tag.Tag,
	}
}

func (s *HttpServer) tagModelsToResponse(tags []query.Tag) []tagResponse {
	responses := make([]tagResponse, len(tags))

	for i, t := range tags {
		responses[i] = tagResponse{
			Id:  t.Id,
			Tag: t.Tag,
		}
	}

	return responses
}
