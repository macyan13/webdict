package query

type SingleTag struct {
	Id       string
	AuthorId string
}

type SingleTagHandler struct {
	tagRepo TagViewRepository
}

func NewSingleTagHandler(tagRepo TagViewRepository) SingleTagHandler {
	return SingleTagHandler{tagRepo: tagRepo}
}

func (h SingleTagHandler) Handle(cmd SingleTag) (TagView, error) {
	return h.tagRepo.GetView(cmd.Id, cmd.AuthorId)
}
