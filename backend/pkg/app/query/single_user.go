package query

// SingleUser get user by ID
type SingleUser struct {
	ID string
}

// SingleUserHandler get tag query handler
type SingleUserHandler struct {
	userRepo   UserViewRepository
	strictSntz *strictSanitizer
}

func NewSingleUserHandler(userRepo UserViewRepository) SingleUserHandler {
	return SingleUserHandler{userRepo: userRepo, strictSntz: newStrictSanitizer()}
}

// Handle performs query to get user by ID
func (h SingleUserHandler) Handle(cmd SingleUser) (UserView, error) {
	view, err := h.userRepo.GetView(cmd.ID)

	if err != nil {
		return UserView{}, err
	}

	view.sanitize(h.strictSntz)
	return view, nil
}
