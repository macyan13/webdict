package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"log"
	"net/http"
	"time"
)

const userContextKey = "user"
const authType = "Bearer"

var ErrInvalidCredentials = errors.New("auth: can not authenticate, invalid email or password")
var ErrExpiredRefreshToken = errors.New("auth: can not refresh auth token, refresh token is expired")

type tokener interface {
	generateToken(email string, expiresAt time.Time) (string, error)
	parseToken(signedToken string) (*JWTClaim, error)
}

type Handler struct {
	userRepo user.Repository
	tokener  tokener
	cipher   Cipher
	params   Params
}

func NewHandler(userRepo user.Repository, cipher Cipher, params Params) *Handler {
	return &Handler{userRepo: userRepo, tokener: jwtTokener{params: params}, cipher: cipher, params: params}
}

func (h Handler) Authenticate(email, password string) (AuthenticationToken, error) {
	usr, err := h.userRepo.GetByEmail(email)

	if err == user.NotFoundErr {
		return AuthenticationToken{}, ErrInvalidCredentials
	}

	if err != nil {
		return AuthenticationToken{}, err
	}

	if !h.cipher.ComparePasswords(usr.Password(), password) {
		return AuthenticationToken{}, ErrInvalidCredentials
	}

	return h.generateAuthToken(email)
}

func (h Handler) GenerateRefreshToken(email string) (RefreshToken, error) {
	expiresAt := time.Now().Add(h.params.RefreshTTL)
	token, err := h.tokener.generateToken(email, expiresAt)

	if err != nil {
		return RefreshToken{}, err
	}

	return RefreshToken{
		Token: token,
	}, nil
}

func (h Handler) Refresh(token string) (AuthenticationToken, error) {
	claims, err := h.tokener.parseToken(token)

	if err != nil {
		return AuthenticationToken{}, err
	}

	return h.generateAuthToken(claims.Email)
}

func (h Handler) generateAuthToken(email string) (AuthenticationToken, error) {
	token, err := h.tokener.generateToken(email, time.Now().Add(h.params.AuthTTL))

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
			return
		}

		claims, err := h.tokener.parseToken(token)

		if err != nil {
			log.Printf("[Error] Can not parse auth token: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		usr, err := h.userRepo.GetByEmail(claims.Email)

		if err == user.NotFoundErr {
			log.Printf("[Error] Attempt to authenticate with not existing user and valid token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Printf("[Error] Can not get user from DB: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
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
		return User{}, fmt.Errorf("can not cast authorised user")
	}

	return usr, nil
}

// "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGVtYWlsLmNvbSIsImV4cCI6MTY3MTEzNjU1OX0.92rJjpLYfmjJIAMcPU4j6gxoxG_ieDpHOP3iBameQXCX8i3h-cLiFrJtJOLNU1Wk_maL1FocJbo5jv2MULuxyw"
func (h Handler) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 8 && headerValue[0:6] == authType {
		return headerValue[7:]
	}

	return ""
}
