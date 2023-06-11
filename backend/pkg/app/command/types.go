package command

// Cipher service to generate password hash before saving a user to DB
type Cipher interface {
	GenerateHash(pwd string) (string, error)
	ComparePasswords(hashedPwd, plainPwd string) bool
}
