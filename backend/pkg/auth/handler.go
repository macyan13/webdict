package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"log"
	"net/http"
	"strings"
	"time"
)

const userContextKey = "user"

var ErrInvalidCredentials = errors.New("auth: can not authenticate, invalid email or password")
var ErrExpiredRefreshToken = errors.New("auth: can not refresh auth token, refresh token is expired")

// todo move to config
var jwtKey = []byte("supersecretkey")
var authType = "Bearer"

type Handler struct {
	userRepo user.Repository
}

func NewHandler(userRepo user.Repository) *Handler {
	return &Handler{userRepo: userRepo}
}

func (h Handler) Authenticate(email, password string) (AuthenticationToken, error) {
	usr := h.userRepo.GetByEmail(email)

	if usr == nil || usr.IsPasswordValid(password) {
		return AuthenticationToken{}, ErrInvalidCredentials
	}

	return h.generateAuthToken(email)
}

func (h Handler) GenerateRefreshToken(email string) (RefreshToken, error) {
	expiresAt := time.Now().Add(time.Hour * 24) // todo: move to configs
	token, err := h.generateToken(email, expiresAt)

	if err != nil {
		return RefreshToken{}, err
	}

	return RefreshToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (h Handler) Refresh(token string) (AuthenticationToken, error) {
	claims, err := h.parseToken(token)

	if err != nil {
		return AuthenticationToken{}, err
	}

	return h.generateAuthToken(claims.Email)
}

func (h Handler) generateToken(email string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h Handler) parseToken(signedToken string) (*JWTClaim, error) {
	claims := JWTClaim{}
	_, err := jwt.ParseWithClaims(
		signedToken,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, ErrExpiredRefreshToken
	}

	return &claims, nil
}

func (h Handler) generateAuthToken(email string) (AuthenticationToken, error) {
	token, err := h.generateToken(email, time.Now().Add(time.Minute*10)) // todo: move to configs

	if err != nil {
		return AuthenticationToken{}, err
	}

	return AuthenticationToken{
		Token: token,
		Type:  authType,
	}, nil
}

func (h Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := h.tokenFromHeader(c.Request)

		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		claims, err := h.parseToken(token)

		if err != nil {
			log.Printf("[Error] Can not parse auth token: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		usr := h.userRepo.GetByEmail(claims.Email)

		if usr != nil {
			log.Printf("[Error] Attempt to authenticate with not existing user and valid token")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set(userContextKey, User{
			Id:    usr.Id(),
			Email: usr.Email(),
			Role:  usr.Role(),
		})
	}
}

func (h Handler) UserFromContext(c *gin.Context) (User, error) {
	value, exists := c.Get(userContextKey)

	if !exists {
		return User{}, fmt.Errorf("can not get authorised user")
	}

	usr, ok := value.(User)

	if !ok {
		return User{}, fmt.Errorf("can not get authorised user")
	}

	return usr, nil
}

func (h Handler) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == authType {
		return headerValue[7:]
	}

	return ""
}
