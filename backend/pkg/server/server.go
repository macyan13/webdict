package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/macyan13/webdict/backend/pkg/repository"
	"log"
	"net/http"
	"runtime"
)

const adminEmail = "admin@email.com"
const adminPassword = "12345678"

type HttpServer struct {
	engine      *gin.Engine
	app         app.Application
	authHandler auth.Handler
}

func InitServer() *HttpServer {
	router := gin.Default()
	// 	"github.com/gin-contrib/cors"
	// router.Use(cors.Default()) - middleware for CORS support, maybe add later
	userRepo := repository.NewUserRepository()
	cipher := auth.Cipher{}

	s := HttpServer{
		engine:      router,
		app:         newApplication(userRepo, cipher),
		authHandler: *auth.NewHandler(userRepo, cipher),
	}

	s.BuildRoutes()
	s.PopulateInitData()
	return &s
}

func newApplication(userRepo user.Repository, cipher auth.Cipher) app.Application {
	tagRepo := repository.NewTagRepository()
	translationRepo := repository.NewTranslationRepository(*tagRepo)

	cmd := app.Commands{
		AddTranslation:    command.NewAddTranslationHandler(translationRepo, tagRepo),
		UpdateTranslation: command.NewUpdateTranslationHandler(translationRepo, tagRepo),
		DeleteTranslation: command.NewDeleteTranslationHandler(translationRepo, tagRepo),
		AddTag:            command.NewAddTagHandler(tagRepo),
		UpdateTag:         command.NewUpdateTagHandler(tagRepo),
		DeleteTag:         command.NewDeleteTagHandler(tagRepo),
		AddUser:           command.NewAddUserHandler(userRepo, cipher),
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

func (s *HttpServer) Run() error {
	err := s.engine.Run(":4000") // todo: move to config

	if err != nil {
		log.Printf("HttpServer - there was an error calling Run on engine: %v", err)
		return err
	}

	return nil
}

func (s *HttpServer) PopulateInitData() {
	// Ignore errors for now
	s.app.Commands.AddUser.Handle(command.AddUser{
		Name:     "admin",
		Email:    adminEmail,
		Password: adminPassword,
	})
}

func (s *HttpServer) unauthorised(c *gin.Context, err error) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.Printf("[Error] Can not authorise action: %s:%s: %v", frame.File, frame.Function, err)
	c.JSON(http.StatusUnauthorized, nil)
}

func (s *HttpServer) badRequest(c *gin.Context, err error) {
	log.Printf("[Error] Can not handle request - %v", err)
	c.JSON(http.StatusBadRequest, nil)
}
