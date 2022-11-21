package server

import (
	"fmt"
)

func (s *HttpServer) BuildRoutes() {
	router := s.engine

	v1 := router.Group("/v1/api")
	{
		translationApi := v1.Group("/translations", s.authHandler.Middleware())
		{
			translationApi.POST("", s.CreateTranslation())
			translationApi.GET("/last", s.GetLastTranslations())
			translationApi.PUT(fmt.Sprintf("/:%s", translationIdParam), s.UpdateTranslation())
			translationApi.GET(fmt.Sprintf("/:%s", translationIdParam), s.GetTranslationById())
			translationApi.DELETE(fmt.Sprintf("/:%s", translationIdParam), s.DeleteTranslationById())
		}

		tagApi := v1.Group("/tags", s.authHandler.Middleware())
		{
			tagApi.POST("", s.CreateTag())
			tagApi.GET("", s.GetTags())
			tagApi.PUT(fmt.Sprintf("/:%s", tagIdParam), s.UpdateTag())
			tagApi.GET(fmt.Sprintf("/:%s", tagIdParam), s.GetTagById())
			tagApi.DELETE(fmt.Sprintf("/:%s", tagIdParam), s.DeleteTagById())
		}

		authApi := v1.Group("/auth")
		{
			authApi.POST("/signin", s.CreateTag())
			authApi.POST("refresh", s.CreateTag())
		}
	}
}
