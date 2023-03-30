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

		id, err := s.app.Commands.AddTranslation.Handle(command.AddTranslation{
			Transcription: request.Transcription,
			Target:        request.Target,
			Text:          request.Source,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorID:      user.ID,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not create new translation: %v", err))
			return
		}

		c.JSON(http.StatusCreated, idResponse{ID: id})
	}
}

func (s *HTTPServer) GetLastTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		var pageSize, page int
		requestPageSize := c.Query("pageSize")
		pageSize, _ = strconv.Atoi(requestPageSize)

		requestPage := c.Query("page")
		page, _ = strconv.Atoi(requestPage)

		lastViews, err := s.app.Queries.LastTranslations.Handle(query.LastTranslations{
			AuthorID: user.ID,
			PageSize: pageSize,
			Page:     page,
			TagIds:   c.QueryArray("tagId"),
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not return last translations, can not perform DB query - %v", err))
			return
		}

		c.JSON(http.StatusOK, lastTranslationsResponse{
			Translations: s.translationViewsToResponse(lastViews.Views),
			TotalPages:   lastViews.TotalPages,
		})
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

		if err := s.app.Commands.UpdateTranslation.Handle(command.UpdateTranslation{
			ID:            c.Param(translationIDParam),
			Target:        request.Target,
			Transcription: request.Transcription,
			Source:        request.Source,
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
		Target:        translation.Target,
		Source:        translation.Source,
		Example:       translation.Example,
		Tags:          tags,
	}
}
