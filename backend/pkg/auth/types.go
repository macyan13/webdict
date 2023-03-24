package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"time"
)

type JWTClaim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthenticationToken struct {
	Token string
	Type  string
}

type RefreshToken struct {
	Token string
}

type User struct {
	ID    string
	Email string
	Role  user.Role
}

func (u User) IsAdmin() bool {
	return u.Role == user.Admin
}

type Params struct {
	AuthTTL    time.Duration
	RefreshTTL time.Duration
	Secret     string
}
