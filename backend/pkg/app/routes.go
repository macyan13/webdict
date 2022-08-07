package app

import "github.com/gin-gonic/gin"

func (server *Server) Routes() *gin.Engine {
	router := server.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		v1.GET("/status", server.ApiStatus())
	}

	return router
}
