package app

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"log"
	"net/http"
)

const translationIdParam = "translationId"

func (s *Server) CreateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translation.Request
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse new Translation request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		request.AuthorId = s.userService.GetByEmail(adminEmail).Id
		if err := s.translationService.CreateTranslation(request); err != nil {
			log.Printf("[Error] Can not create new Translation: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *Server) GetTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, s.translationService.GetTranslations())
	}
}

func (s *Server) UpdateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translation.Request

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse new Translation request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		request.AuthorId = s.userService.GetByEmail(adminEmail).Id

		if err := s.translationService.UpdateTranslation(c.Param(translationIdParam), request); err != nil {
			log.Printf("[Error] Can not Update Existing Translation: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		response := map[string]any{
			"status": "success",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetTranslationById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		record := s.translationService.GetById(c.Param(translationIdParam))

		if record == nil {
			c.JSON(http.StatusBadRequest, "Can not find requested record")
			return
		}

		c.JSON(http.StatusOK, record)
	}
}

func (s *Server) DeleteTranslationById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		if err := s.translationService.DeleteById(c.Param(translationIdParam)); err != nil {
			log.Printf("[Error] Can not delete translation: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		response := map[string]any{
			"status": "success",
		}

		c.JSON(http.StatusOK, response)
	}
}
