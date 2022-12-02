package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type jwtTokener struct {
}

func (t jwtTokener) generateToken(email string, expiresAt time.Time) (string, error) {
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

func (t jwtTokener) parseToken(signedToken string) (*JWTClaim, error) {
	claims := JWTClaim{}
	_, err := jwt.ParseWithClaims(
		signedToken,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)

	return &claims, err
}
