package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		translationApi := v1.Group("/translations")
		{
			translationApi.POST("", s.CreateTranslation())
			translationApi.GET("", s.GetTranslations())
			translationApi.PUT(fmt.Sprintf("/:%s", translationIdParam), s.UpdateTranslation())
			translationApi.GET(fmt.Sprintf("/:%s", translationIdParam), s.GetTranslationById())
			translationApi.DELETE(fmt.Sprintf("/:%s", translationIdParam), s.DeleteTranslationById())
		}
	}

	return router
}
