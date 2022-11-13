package query

type AllTags struct {
	AuthorId string
}

// todo: maybe move query repositories to one definition
type AllTagsRepository interface {
	GetAll(authorId string) []Tag
}

type AllTagsHandler struct {
	tagRepo AllTagsRepository
}

func NewAllTagsHandler(tagRepository AllTagsRepository) AllTagsHandler {
	return AllTagsHandler{tagRepo: tagRepository}
}

func (h AllTagsHandler) Handle(cmd AllTags) []Tag {
	return h.tagRepo.GetAll(cmd.AuthorId)
}
