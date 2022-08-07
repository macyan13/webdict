package app

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	router *gin.Engine
}

func NewServer(router *gin.Engine) *Server {
	return &Server{
		router: router,
	}
}

func (server *Server) Run() error {
	routes := server.Routes()

	err := routes.Run(":4000") // todo: move to config

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
