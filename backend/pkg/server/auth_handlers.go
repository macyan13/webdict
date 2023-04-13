package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/auth"
	"log"
	"net/http"
	"time"
)

const refreshTokenCookieName = "refreshToken"

func (s *HTTPServer) SighIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var request SignInRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			s.badRequest(c, fmt.Errorf("[Error] Can not parse SighIn request: %v", err))
			return
		}

		authToken, err := s.authHandler.Authenticate(request.Email, request.Password)

		if err != nil {
			if err != auth.ErrInvalidCredentials {
				log.Printf("[Error] Can not handle auth request: %v", err)
			}
			c.JSON(http.StatusUnauthorized, nil)
			return
		}

		refreshToken, err := s.authHandler.GenerateRefreshToken(request.Email)

		if err != nil {
			s.unauthorized(c, fmt.Errorf("[Error] Can not generate Refresh token: %v", err))
			return
		}

		c.SetCookie(refreshTokenCookieName, refreshToken.Token, int(time.Now().Add(s.opts.Auth.TTL.Cookie).Unix()), "/", s.opts.WebdictURL, false, true)

		c.JSON(http.StatusOK, AuthTokenResponse{
			AccessToken: authToken.Token,
			Type:        authToken.Type,
		})
	}
}

func (s *HTTPServer) Refresh() gin.HandlerFunc {
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
