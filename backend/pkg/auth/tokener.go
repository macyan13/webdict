package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type jwtTokener struct {
	params Params
}

func (t jwtTokener) generateToken(email string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	tokenString, err := token.SignedString([]byte(t.params.Secret))
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
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(t.params.Secret), nil
		},
	)

	return &claims, err
}
