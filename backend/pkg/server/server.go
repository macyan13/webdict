package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/macyan13/webdict/backend/pkg/repository"
	"log"
)

const adminEmail = "admin@email.com"

type HttpServer struct {
	engine   *gin.Engine
	app      app.Application
	userRepo user.Repository // todo: remove after auth implementation
}

func InitServer() *HttpServer {
	router := gin.Default()
	// 	"github.com/gin-contrib/cors"
	// router.Use(cors.Default()) - middleware for CORS support, maybe add later
	userRepo := repository.NewUserRepository()

	s := HttpServer{
		engine:   router,
		app:      newApplication(userRepo),
		userRepo: userRepo,
	}

	s.BuildRoutes()
	s.PopulateInitData()
	return &s
}

func newApplication(userRepo user.Repository) app.Application {
	tagRepo := repository.NewTagRepository()
	translationRepo := repository.NewTranslationRepository(*tagRepo)

	cmd := app.Commands{
		AddTranslation:    command.NewAddTranslationHandler(translationRepo, tagRepo),
		UpdateTranslation: command.NewUpdateTranslationHandler(translationRepo, tagRepo),
		DeleteTranslation: command.NewDeleteTranslationHandler(translationRepo, tagRepo),
		AddTag:            command.NewAddTagHandler(tagRepo),
		UpdateTag:         command.NewUpdateTagHandler(tagRepo),
		DeleteTag:         command.NewDeleteTagHandler(tagRepo),
		AddUser:           command.NewAddUserHandler(userRepo),
	}

	queries := app.Queries{
		SingleTranslation: query.NewSingleTranslationHandler(translationRepo),
		LastTranslations:  query.NewLastTranslationsHandler(translationRepo),
		SingleTag:         query.NewSingleTagHandler(tagRepo),
		AllTags:           query.NewAllTagsHandler(tagRepo),
	}

	return app.Application{
		Commands: cmd,
		Queries:  queries,
	}
}

func NewAppServer(router *gin.Engine, app app.Application, userRepo user.Repository) *HttpServer {
	return &HttpServer{
		engine:   router,
		app:      app,
		userRepo: userRepo,
	}
}

func (s *HttpServer) Run() error {
	err := s.engine.Run(":4000") // todo: move to config

	if err != nil {
		log.Printf("HttpServer - there was an error calling Run on engine: %v", err)
		return err
	}

	return nil
}

func (s *HttpServer) PopulateInitData() {
	// Ignore errors for now "admin", adminEmail, "password"
	s.app.Commands.AddUser.Handle(command.AddUser{
		Name:     "admin",
		Email:    adminEmail,
		Password: "password",
	})
}
