package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
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
	Token     string
	ExpiresAt time.Time
}

type User struct {
	Id    string
	Email string
	Role  user.Role
}
