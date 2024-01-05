package query

import "github.com/go-playground/validator/v10"

// SingleUser get user by ID
type SingleUser struct {
	ID string `validate:"required"`
}

// SingleUserHandler get tag query handler
type SingleUserHandler struct {
	userRepo   UserViewRepository
	validator  *validator.Validate
	strictSntz *strictSanitizer
}

func NewSingleUserHandler(userRepo UserViewRepository, validate *validator.Validate) SingleUserHandler {
	return SingleUserHandler{userRepo: userRepo, validator: validate, strictSntz: newStrictSanitizer()}
}

// Handle performs query to get user by ID
func (h SingleUserHandler) Handle(cmd SingleUser) (UserView, error) {
	if err := h.validator.Struct(cmd); err != nil {
		return UserView{}, err
	}

	view, err := h.userRepo.GetView(cmd.ID)

	if err != nil {
		return UserView{}, err
	}

	view.sanitize(h.strictSntz)
	return view, nil
}
