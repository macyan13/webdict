package server

import (
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"log"
	"net/http"
	"time"
)

const refreshTokenCookieName = "refreshToken"

func (s *HttpServer) SighIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request SignInRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			log.Printf("[Error] Can not parse SighIn request: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		authToken, err := s.authHandler.Authenticate(request.Email, request.Password)

		if err != nil {
			if err != auth.ErrInvalidCredentials {
				log.Printf("[Error] Can not handle auth request: %v", err)
			}
			c.JSON(http.StatusUnauthorized, nil)
		}

		refreshToken, err := s.authHandler.GenerateRefreshToken(request.Email)

		if err != nil {
			log.Printf("[Error] Can not generate Refresh token: %v", err)
			c.JSON(http.StatusUnauthorized, nil)
		}

		c.SetCookie(refreshTokenCookieName, refreshToken.Token, int(time.Now().Add(time.Hour*24*30*12).Unix()), "/", "", false, true) // todo: get domain from config

		c.JSON(http.StatusOK, AuthTokenResponse{
			AccessToken: authToken.Token,
			Type:        authToken.Type,
		})
	}
}

func (s *HttpServer) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		refreshToken, err := c.Cookie(refreshTokenCookieName)

		if err != nil {
			c.JSON(http.StatusBadRequest, nil)
		}

		authToken, err := s.authHandler.Refresh(refreshToken)

		if err != nil {
			if err != auth.ErrExpiredRefreshToken {
				log.Printf("[Error] Can not handle Refresh token request: %v", err)
			}
			c.JSON(http.StatusUnauthorized, nil)
		}

		c.JSON(http.StatusOK, AuthTokenResponse{
			AccessToken: authToken.Token,
			Type:        authToken.Type,
		})
	}
}
