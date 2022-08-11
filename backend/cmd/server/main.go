package main

import (
	"fmt"
	"github.com/Yan-Matskevich/webdict/backend/pkg/app"
	"github.com/Yan-Matskevich/webdict/backend/pkg/domain/service"
	"github.com/Yan-Matskevich/webdict/backend/pkg/repository"
	"github.com/gin-gonic/gin"
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

	translationService := service.NewTranslationService(repository.NewTranslationRepository())
	return app.NewServer(router, *translationService)
}
