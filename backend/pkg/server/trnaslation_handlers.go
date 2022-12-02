package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
	"strconv"
)

const translationIdParam = "translationId"

func (s *HttpServer) CreateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse new translation request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorised(c, err)
		}

		if err := s.app.Commands.AddTranslation.Handle(command.AddTranslation{
			Transcription: request.Transcription,
			Translation:   request.Translation,
			Text:          request.Text,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorId:      user.Id,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not create new translation: %v", err))
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *HttpServer) GetLastTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorised(c, err)
		}

		limit, err := strconv.Atoi(c.Query("limit"))

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not return last translations, can not parse limit param: %v", c.Query("limit")))
			return
		}

		translations := s.translationModelsToResponse(s.app.Queries.LastTranslations.Handle(query.LastTranslations{
			AuthorId: user.Id,
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
			s.badRequest(c, fmt.Errorf("can not parse update translation request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorised(c, err)
		}

		if err := s.app.Commands.UpdateTranslation.Handle(command.UpdateTranslation{
			Id:            c.Param(translationIdParam),
			Transcription: request.Transcription,
			Translation:   request.Translation,
			Text:          request.Text,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorId:      user.Id,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not Update Existing translation: %v", err))
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

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorised(c, err)
		}

		translation := s.app.Queries.SingleTranslation.Handle(query.SingleTranslation{
			Id:       c.Param(translationIdParam),
			AuthorId: user.Id,
		})

		if translation == nil {
			s.badRequest(c, fmt.Errorf("can not find requested record"))
			return
		}

		c.JSON(http.StatusOK, s.translationModelToResponse(*translation))
	}
}

func (s *HttpServer) DeleteTranslationById() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorised(c, err)
		}

		if err := s.app.Commands.DeleteTranslation.Handle(command.DeleteTranslation{
			Id:       c.Param(translationIdParam),
			AuthorId: user.Id,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not delete translation: %v", err))
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
