package query

// SingleTag get tag by ID and authorID query
type SingleTag struct {
	ID       string
	AuthorID string
}

// SingleTagHandler get tag query handler
type SingleTagHandler struct {
	tagRepo    TagViewRepository
	strictSntz *strictSanitizer
}

func NewSingleTagHandler(tagRepo TagViewRepository) SingleTagHandler {
	return SingleTagHandler{tagRepo: tagRepo, strictSntz: newStrictSanitizer()}
}

// Handle performs query to get tag by ID and authorID
func (h SingleTagHandler) Handle(cmd SingleTag) (TagView, error) {
	view, err := h.tagRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return TagView{}, err
	}

	view.sanitize(h.strictSntz)
	return view, nil
}
