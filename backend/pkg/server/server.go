package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/app"
	"github.com/macyan13/webdict/backend/pkg/app/command"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"log"
)

const adminEmail = "admin@email.com"

type HttpServer struct {
	engine             *gin.Engine
	translationService translation.Service
	tagService         tag.Service
	userService        user.Service
	app                app.Application
}

func NewAppServer(router *gin.Engine, translationService translation.Service, tagService tag.Service, userService user.Service) *HttpServer {
	return &HttpServer{
		engine:             router,
		translationService: translationService,
		tagService:         tagService,
		userService:        userService,
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
