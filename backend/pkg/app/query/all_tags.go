package query

type AllTags struct {
	AuthorID string
}

type AllTagsHandler struct {
	tagRepo TagViewRepository
}

func NewAllTagsHandler(tagRepository TagViewRepository) AllTagsHandler {
	return AllTagsHandler{tagRepo: tagRepository}
}

func (h AllTagsHandler) Handle(cmd AllTags) ([]TagView, error) {
	return h.tagRepo.GetAllViews(cmd.AuthorID)
}
