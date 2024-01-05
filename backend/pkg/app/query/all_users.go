package query

// AllUsersHandler get all users
type AllUsersHandler struct {
	userRepo  UserViewRepository
	sanitizer *strictSanitizer
}

func NewAllUsersHandler(userRepo UserViewRepository) AllUsersHandler {
	return AllUsersHandler{userRepo: userRepo, sanitizer: newStrictSanitizer()}
}

// Handle performs query to receive all tags for author
func (h AllUsersHandler) Handle() ([]UserView, error) {
	users, err := h.userRepo.GetAllViews()

	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].sanitize(h.sanitizer)
	}

	return users, nil
}
