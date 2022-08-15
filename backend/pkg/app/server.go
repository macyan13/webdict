package app

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/service"
	"log"
)

type Server struct {
	router             *gin.Engine
	translationService service.TranslationService
}

func NewServer(router *gin.Engine, translationService service.TranslationService) *Server {
	return &Server{
		router:             router,
		translationService: translationService,
	}
}

func (s *Server) Run() error {
	routes := s.Routes()

	err := routes.Run(":4000") // todo: move to config

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
