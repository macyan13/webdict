package app

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/macyan13/webdict/backend/pkg/repository"
)

func initTestServer() *Server {
	router := gin.Default()
	tagRepository := repository.NewTagRepository()
	userRepository := repository.NewUserRepository()
	s := NewServer(
		router,
		*translation.NewService(repository.NewTranslationRepository(), tagRepository, userRepository),
		*tag.NewService(tagRepository),
		*user.NewService(userRepository),
	)
	s.BuildRoutes()
	s.PopulateInitData()
	return s
}
