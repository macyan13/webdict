package query

import "github.com/go-playground/validator/v10"

// SingleTag get tag by ID and authorID query
type SingleTag struct {
	ID       string `validate:"required"`
	AuthorID string `validate:"required"`
}

// SingleTagHandler get tag query handler
type SingleTagHandler struct {
	tagRepo    TagViewRepository
	validator  *validator.Validate
	strictSntz *strictSanitizer
}

func NewSingleTagHandler(tagRepo TagViewRepository, validate *validator.Validate) SingleTagHandler {
	return SingleTagHandler{tagRepo: tagRepo, validator: validate, strictSntz: newStrictSanitizer()}
}

// Handle performs query to get tag by ID and authorID
func (h SingleTagHandler) Handle(cmd SingleTag) (TagView, error) {
	if err := h.validator.Struct(cmd); err != nil {
		return TagView{}, err
	}

	view, err := h.tagRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return TagView{}, err
	}

	view.sanitize(h.strictSntz)
	return view, nil
}
