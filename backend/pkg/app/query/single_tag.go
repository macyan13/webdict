package query

type SingleTag struct {
	Id       string
	AuthorId string
}

type SingleTagRepository interface {
	getTag(id, authorId string) *Tag
}

type SingleTagHandler struct {
	tagRepo SingleTagRepository
}

func NewSingleTagHandler(tagRepo SingleTagRepository) SingleTagHandler {
	return SingleTagHandler{tagRepo: tagRepo}
}

func (h SingleTagHandler) Handle(cmd SingleTag) *Tag {
	return h.tagRepo.getTag(cmd.Id, cmd.AuthorId)
}
