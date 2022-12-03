package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCipher_GenerateHash(t *testing.T) {
	cipher := Cipher{}
	generatedHash1, err := cipher.GenerateHash("testPassword")

	assert.Nil(t, err)
	assert.True(t, cipher.ComparePasswords(generatedHash1, "testPassword"))
	assert.False(t, cipher.ComparePasswords(generatedHash1, "invalidTestPassword"))
}
