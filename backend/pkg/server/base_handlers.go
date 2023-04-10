package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *HTTPServer) ServeStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	}
}
