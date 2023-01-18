package query

// AllTags get all tags for author query
type AllTags struct {
	AuthorID string
}

// AllTagsHandler get all tags for author query
type AllTagsHandler struct {
	tagRepo TagViewRepository
}

func NewAllTagsHandler(tagRepository TagViewRepository) AllTagsHandler {
	return AllTagsHandler{tagRepo: tagRepository}
}

// Handle performs query to receive all tags for author
func (h AllTagsHandler) Handle(query AllTags) ([]TagView, error) {
	tags, err := h.tagRepo.GetAllViews(query.AuthorID)

	if err != nil {
		return nil, err
	}

	for i := range tags {
		tags[i].sanitize()
	}

	return tags, nil
}
