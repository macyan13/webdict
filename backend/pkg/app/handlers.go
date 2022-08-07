package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (server *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status": "success",
			"data":   "Web dictionary API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}
