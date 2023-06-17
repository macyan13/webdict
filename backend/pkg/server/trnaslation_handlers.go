package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
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
			Source:        request.Source,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorID:      user.ID,
			LangID:        request.LangID,
		})

		if err != nil {
			if err == translation.ErrSourceAlreadyExists {
				s.badRequest(c, fmt.Errorf("translation with source %s already exists", request.Source))
				return
			}
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
		pageSize, _ = strconv.Atoi(c.Query("pageSize"))
		page, _ = strconv.Atoi(c.Query("page"))

		lastViews, err := s.app.Queries.LastTranslations.Handle(query.LastTranslations{
			AuthorID: user.ID,
			PageSize: pageSize,
			Page:     page,
			TagIds:   c.QueryArray("tagId[]"),
			LangID:   c.Query("langId"),
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not return last translations - %v", err))
			return
		}

		c.JSON(http.StatusOK, lastTranslationsResponse{
			Translations: s.translationViewsToResponse(lastViews.Views),
			TotalRecords: lastViews.TotalRecords,
		})
	}
}

func (s *HTTPServer) GetRandomTranslations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		user, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
		}

		limit, _ := strconv.Atoi(c.Query("limit"))

		lastViews, err := s.app.Queries.RandomTranslations.Handle(query.RandomTranslations{
			AuthorID: user.ID,
			LangID:   c.Query("langId"),
			TagIds:   c.QueryArray("tagId[]"),
			Limit:    limit,
		})

		if err != nil {
			s.badRequest(c, fmt.Errorf("can not return random translations - %v", err))
			return
		}

		c.JSON(http.StatusOK, randomTranslationsResponse{
			Translations: s.translationViewsToResponse(lastViews.Views),
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

		if err = s.app.Commands.UpdateTranslation.Handle(command.UpdateTranslation{
			ID:            c.Param(translationIDParam),
			Target:        request.Target,
			Transcription: request.Transcription,
			Source:        request.Source,
			Example:       request.Example,
			TagIds:        request.TagIds,
			AuthorID:      user.ID,
			LangID:        request.LangID,
		}); err != nil {
			if err == translation.ErrSourceAlreadyExists {
				s.badRequest(c, fmt.Errorf("translation with source %s already exists", request.Source))
				return
			}
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

		if err = s.app.Commands.DeleteTranslation.Handle(command.DeleteTranslation{
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

func (s *HTTPServer) translationViewToResponse(view query.TranslationView) translationResponse {
	tags := make([]tagResponse, len(view.Tags))

	for i, tag := range view.Tags {
		tags[i] = tagResponse{
			ID:  tag.ID,
			Tag: tag.Tag,
		}
	}

	return translationResponse{
		ID:            view.ID,
		CreatedAt:     view.CreatedAd,
		Transcription: view.Transcription,
		Target:        view.Target,
		Source:        view.Source,
		Example:       view.Example,
		Tags:          tags,
		Lang: langResponse{
			ID:   view.Lang.ID,
			Name: view.Lang.Name,
		},
	}
}
