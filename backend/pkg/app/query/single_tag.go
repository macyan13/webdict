package query

type SingleTag struct {
	ID       string
	AuthorID string
}

type SingleTagHandler struct {
	tagRepo TagViewRepository
}

func NewSingleTagHandler(tagRepo TagViewRepository) SingleTagHandler {
	return SingleTagHandler{tagRepo: tagRepo}
}

func (h SingleTagHandler) Handle(cmd SingleTag) (TagView, error) {
	return h.tagRepo.GetView(cmd.ID, cmd.AuthorID)
}
