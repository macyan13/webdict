package server

import (
	"fmt"
)

func (s *HTTPServer) buildRoutes() {
	router := s.engine

	router.Group("/").GET("", s.ServeStatic())

	v1 := router.Group("/v1/api")
	{
		authAPI := v1.Group("/auth")
		authAPI.POST("/signin", s.SighIn())
		authAPI.POST("/refresh", s.Refresh())

		translationAPI := v1.Group("/translations", s.authHandler.Middleware())
		translationAPI.POST("", s.CreateTranslation())
		translationAPI.GET("/last", s.GetLastTranslations())
		translationAPI.GET("/random", s.GetRandomTranslations())
		translationAPI.PUT(fmt.Sprintf("/:%s", translationIDParam), s.UpdateTranslation())
		translationAPI.GET(fmt.Sprintf("/:%s", translationIDParam), s.GetTranslationByID())
		translationAPI.DELETE(fmt.Sprintf("/:%s", translationIDParam), s.DeleteTranslationByID())

		tagAPI := v1.Group("/tags", s.authHandler.Middleware())
		tagAPI.POST("", s.CreateTag())
		tagAPI.GET("", s.GetTags())
		tagAPI.PUT(fmt.Sprintf("/:%s", tagIDParam), s.UpdateTag())
		tagAPI.GET(fmt.Sprintf("/:%s", tagIDParam), s.GetTagByID())
		tagAPI.DELETE(fmt.Sprintf("/:%s", tagIDParam), s.DeleteTagByID())

		userAPI := v1.Group("/users", s.authHandler.Middleware(), s.authHandler.AdminMiddleware())
		userAPI.POST("", s.CreateUser())
		userAPI.PUT(fmt.Sprintf("/:%s", userIDParam), s.UpdateUser())
		userAPI.GET("", s.GetUsers())
		userAPI.GET(fmt.Sprintf("/:%s", userIDParam), s.GetUserByID())

		roleAPI := v1.Group("/roles", s.authHandler.Middleware(), s.authHandler.AdminMiddleware())
		roleAPI.GET("", s.GetRoles())

		langAPI := v1.Group("/langs", s.authHandler.Middleware())
		langAPI.GET("", s.GetLangs())
		langAPI.POST("", s.CreateLang())
		langAPI.PUT(fmt.Sprintf("/:%s", langIDParam), s.UpdateLang())
		langAPI.GET(fmt.Sprintf("/:%s", langIDParam), s.GetLangByID())
		langAPI.DELETE(fmt.Sprintf("/:%s", langIDParam), s.DeleteLangByID())

		profileAPI := v1.Group("/profile", s.authHandler.Middleware())
		profileAPI.GET("", s.GetProfile())
		profileAPI.PUT("", s.UpdateProfile())
	}
}
