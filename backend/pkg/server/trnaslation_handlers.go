package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"log"
	"net/http"
	"strconv"
)

const translationIdParam = "translationId"

func (s *HttpServer) CreateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse new Translation request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id

		if err := s.app.Commands.AddTranslation.Handle(command.AddTranslation{
			Transcription: request.Transcription,
			Translation:   request.Translation,
			Text:          request.Text,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorId:      authorId,
		}); err != nil {
			log.Printf("[Error] Can not create new Translation: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *HttpServer) GetLastTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		authorId := s.userService.GetByEmail(adminEmail).Id
		limit, err := strconv.Atoi(c.Query("limit"))

		if err != nil {
			log.Printf("[Error] Can not return last translations, can not parse limit param: %v", c.Query("limit"))
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		translations := s.translationModelsToResponse(s.app.Queries.LastTranslations.Handle(query.LastTranslations{
			AuthorId: authorId,
			Limit:    limit,
		}))

		c.JSON(http.StatusOK, translations)
	}
}

func (s *HttpServer) UpdateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translationRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse new Translation request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		authorId := s.userService.GetByEmail(adminEmail).Id

		if err := s.app.Commands.UpdateTranslation.Handle(command.UpdateTranslation{
			Id:            c.Param(translationIdParam),
			Transcription: request.Transcription,
			Translation:   request.Translation,
			Text:          request.Text,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorId:      authorId,
		}); err != nil {
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

func (s *HttpServer) GetTranslationById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// todo get user details from auth context
		authorId := s.userService.GetByEmail(adminEmail).Id
		translation := s.app.Queries.SingleTranslation.Handle(query.SingleTranslation{
			Id:       c.Param(translationIdParam),
			AuthorId: authorId,
		})

		if translation == nil {
			c.JSON(http.StatusBadRequest, "Can not find requested record")
			return
		}

		c.JSON(http.StatusOK, s.translationModelToResponse(*translation))
	}
}

func (s *HttpServer) DeleteTranslationById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		authorId := s.userService.GetByEmail(adminEmail).Id

		if err := s.app.Commands.DeleteTranslation.Handle(command.DeleteTranslation{
			Id:       c.Param(translationIdParam),
			AuthorId: authorId,
		}); err != nil {
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

func (s *HttpServer) translationModelsToResponse(translations []query.Translation) []translationResponse {
	responses := make([]translationResponse, len(translations))

	for i, translation := range translations {
		responses[i] = s.translationModelToResponse(translation)
	}

	return responses
}

func (s *HttpServer) translationModelToResponse(translation query.Translation) translationResponse {
	tags := make([]tagResponse, len(translation.Tags))

	for i, tag := range translation.Tags {
		tags[i] = tagResponse{
			Id:  tag.Id,
			Tag: tag.Tag,
		}
	}

	return translationResponse{
		Id:            translation.Id,
		CreatedAt:     translation.CreatedAd,
		Transcription: translation.Transcription,
		Translation:   translation.Translation,
		Text:          translation.Text,
		Example:       translation.Example,
		Tags:          tags,
	}
}
