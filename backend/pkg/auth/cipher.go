package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Cipher struct {
}

func (c Cipher) GenerateHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 8)
	return string(hash), err
}

func (c Cipher) ComparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
