package query

type AllTags struct {
	AuthorId string
}

// todo: maybe move query repositories to one definition
type AllTagsRepository interface {
	getAll(authorId string) []Tag
}

type AllTagsHandler struct {
	tagRepo AllTagsRepository
}

func NewAllTagHandlerHandler(tagRepository AllTagsRepository) AllTagsHandler {
	return AllTagsHandler{tagRepo: tagRepository}
}

func (h AllTagsHandler) Handle(cmd AllTags) []Tag {
	return h.tagRepo.getAll(cmd.AuthorId)
}
