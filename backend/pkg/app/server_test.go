package app

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/repository"
)

func initTestServer() *Server {
	router := gin.Default()
	tagRepo := repository.NewTagRepository()
	s := NewServer(
		router,
		*translation.NewService(repository.NewTranslationRepository(), tagRepo),
		*tag.NewService(tagRepo),
	)
	s.BuildRoutes()
	return s
}
