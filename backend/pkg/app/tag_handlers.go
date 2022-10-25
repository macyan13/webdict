package app

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"log"
	"net/http"
)

const tagIdParam = "tagId"

func (s *Server) CreateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request tag.Request
		err := c.ShouldBindJSON(&request)

		if err != nil {
			log.Printf("[Error] Can not parse new Tag request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.tagService.CreateTag(request)

		if err != nil {
			log.Printf("[Error] Can not create new Tag: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *Server) GetTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, s.tagService.GetTags())
	}
}

func (s *Server) UpdateTag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request tag.Request
		err := c.ShouldBindJSON(&request)

		if err != nil {
			log.Printf("[Error] Can not parse new Tag request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		err = s.tagService.UpdateTag(c.Param(tagIdParam), request)

		if err != nil {
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

func (s *Server) GetTagById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		record := s.tagService.GetById(c.Param(tagIdParam))

		if record == nil {
			c.JSON(http.StatusBadRequest, "Can not find requested record")
		}

		c.JSON(http.StatusOK, record)
	}
}

func (s *Server) DeleteTagById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		err := s.tagService.DeleteById(c.Param(tagIdParam))

		if err != nil {
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
