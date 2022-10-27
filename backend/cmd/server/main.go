package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/macyan13/webdict/backend/pkg/repository"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	server := initServer()
	err := server.Run()

	if err != nil {
		return err
	}

	return nil
}

func initServer() *app.Server {
	router := gin.Default()
	// 	"github.com/gin-contrib/cors"
	// router.Use(cors.Default()) - middleware for CORS support, maybe add later

	tagRepository := repository.NewTagRepository()
	userRepository := repository.NewUserRepository()
	s := app.NewServer(
		router,
		*translation.NewService(repository.NewTranslationRepository(), tagRepository, userRepository),
		*tag.NewService(tagRepository),
		*user.NewService(userRepository),
	)
	s.BuildRoutes()
	s.PopulateInitData()
	return s
}
