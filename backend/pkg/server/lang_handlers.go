package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *HTTPServer) GetLangs() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		_, err := s.authHandler.UserFromContext(c)
		if err != nil {
			s.unauthorized(c, err)
			return
		}

		c.JSON(http.StatusOK, s.app.Queries.SupportedLangs.Handle())
	}
}
