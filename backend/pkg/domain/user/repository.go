package user

type Repository interface {
	Exist(email string) bool
	Save(user *User) error
	GetByEmail(email string) *User
}
