package user

type Repository interface {
	Exist(email string) bool
	Save(user *User) error
	GetByEmail(email string) *User
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateUser(name, email, password string) error {
	user := NewUser(name, email, password)
	return s.repository.Save(user)
}

func (s *Service) GetByEmail(email string) *User {
	return s.repository.GetByEmail(email)
}
