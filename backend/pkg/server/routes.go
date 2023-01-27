package server

import (
	"fmt"
)

func (s *HTTPServer) BuildRoutes() {
	router := s.engine

	v1 := router.Group("/v1/api")
	{
		authAPI := v1.Group("/auth")
		authAPI.POST("/signin", s.SighIn())
		authAPI.POST("/refresh", s.Refresh())

		translationAPI := v1.Group("/translations", s.authHandler.Middleware())
		translationAPI.POST("", s.CreateTranslation())
		translationAPI.GET("/last", s.GetLastTranslations())
		translationAPI.PUT(fmt.Sprintf("/:%s", translationIDParam), s.UpdateTranslation())
		translationAPI.GET(fmt.Sprintf("/:%s", translationIDParam), s.GetTranslationByID())
		translationAPI.DELETE(fmt.Sprintf("/:%s", translationIDParam), s.DeleteTranslationByID())

		tagAPI := v1.Group("/tags", s.authHandler.Middleware())
		tagAPI.POST("", s.CreateTag())
		tagAPI.GET("", s.GetTags())
		tagAPI.PUT(fmt.Sprintf("/:%s", tagIDParam), s.UpdateTag())
		tagAPI.GET(fmt.Sprintf("/:%s", tagIDParam), s.GetTagByID())
		tagAPI.DELETE(fmt.Sprintf("/:%s", tagIDParam), s.DeleteTagByID())

		userAPI := v1.Group("/users", s.authHandler.Middleware())
		userAPI.POST("", s.CreateUser())
		userAPI.PUT(fmt.Sprintf("/:%s", userIDParam), s.UpdateUser())
		userAPI.GET("", s.GetUsers())
		userAPI.GET(fmt.Sprintf("/:%s", userIDParam), s.GetUserByID())
	}
}
