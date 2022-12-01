package auth

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestJwtToken_parseExpiredToken(t *testing.T) {
	tokener := jwtTokener{}
	token, err := tokener.generateToken("testEmail", time.Now().Add(-time.Minute))
	assert.Nil(t, err)

	_, err = tokener.parseToken(token)
	assert.True(t, strings.Contains(err.Error(), "token is expired by"))
}

func TestJwtToken_parseGeneratedToken(t *testing.T) {
	tokener := jwtTokener{}
	email := "testEmail"
	token, err := tokener.generateToken(email, time.Now().Add(time.Minute))
	assert.Nil(t, err)

	claims, err := tokener.parseToken(token)
	assert.Nil(t, err)
	assert.Equalf(t, email, claims.Email, "jwtTokener::parseToken - email from generated token (%s) parsed correctly", token)
}
