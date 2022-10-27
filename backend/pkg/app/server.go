package app

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"log"
)

const adminEmail = "admin@email.com"

type Server struct {
	router             *gin.Engine
	translationService translation.Service
	tagService         tag.Service
	userService        user.Service
}

func NewServer(router *gin.Engine, translationService translation.Service, tagService tag.Service, userService user.Service) *Server {
	return &Server{
		router:             router,
		translationService: translationService,
		tagService:         tagService,
		userService:        userService,
	}
}

func (s *Server) Run() error {
	err := s.router.Run(":4000") // todo: move to config

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}

func (s *Server) PopulateInitData() {
	// Ignore errors for now
	s.userService.CreateUser("admin", adminEmail, "password")
}
