package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"net/http"
	"strconv"
)

const translationIDParam = "translationId"

func (s *HTTPServer) CreateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translationRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse new translation request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if err := s.app.Commands.AddTranslation.Handle(command.AddTranslation{
			Transcription: request.Transcription,
			Translation:   request.Translation,
			Text:          request.Text,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorID:      user.ID,
		}); err != nil {
			s.badRequest(c, fmt.Errorf("can not create new translation: %v", err))
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (s *HTTPServer) GetLastTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		var limit int
		requestLimit, exist := c.GetQuery("limit")

		if exist {
			limit, err = strconv.Atoi(requestLimit)

			if err != nil {
				s.badRequest(c, fmt.Errorf("can not return last translations, can not parse limit param: %v, %v", c.Query("limit"), err))
				return
			}
		}

		views, err := s.app.Queries.LastTranslations.Handle(query.LastTranslations{
			AuthorID: user.ID,
			Limit:    limit,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not return last translations, can not perform DB query - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.translationViewsToResponse(views))
	}
}

func (s *HTTPServer) UpdateTranslation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		var request translationRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("can not parse update translation request: %v", err))
			return
		}

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if err := s.app.Commands.UpdateTranslation.Handle(&command.UpdateTranslation{
			ID:            c.Param(translationIDParam),
			Transcription: request.Transcription,
			Translation:   request.Translation,
			Text:          request.Text,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorID:      user.ID,
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

func (s *HTTPServer) GetTranslationByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		view, err := s.app.Queries.SingleTranslation.Handle(query.SingleTranslation{
			ID:       c.Param(translationIDParam),
			AuthorID: user.ID,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not get requested record - %v", err))
			return
		}

		c.JSON(http.StatusOK, s.translationViewToResponse(view))
	}
}

func (s *HTTPServer) DeleteTranslationByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		if err := s.app.Commands.DeleteTranslation.Handle(command.DeleteTranslation{
			ID:       c.Param(translationIDParam),
			AuthorID: user.ID,
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

func (s *HTTPServer) translationViewsToResponse(translations []query.TranslationView) []translationResponse {
	responses := make([]translationResponse, len(translations))

	for i := range translations {
		responses[i] = s.translationViewToResponse(translations[i])
	}

	return responses
}

func (s *HTTPServer) translationViewToResponse(translation query.TranslationView) translationResponse {
	tags := make([]tagResponse, len(translation.Tags))

	for i, tag := range translation.Tags {
		tags[i] = tagResponse{
			ID:  tag.ID,
			Tag: tag.Tag,
		}
	}

	return translationResponse{
		ID:            translation.ID,
		CreatedAt:     translation.CreatedAd,
		Transcription: translation.Transcription,
		Translation:   translation.Translation,
		Text:          translation.Text,
		Example:       translation.Example,
		Tags:          tags,
	}
}
