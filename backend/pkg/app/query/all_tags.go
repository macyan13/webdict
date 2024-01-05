package query

import "github.com/go-playground/validator/v10"

// AllTags get all tags for author query
type AllTags struct {
	AuthorID string `validate:"required"`
}

// AllTagsHandler get all tags for author query
type AllTagsHandler struct {
	tagRepo   TagViewRepository
	sanitizer *strictSanitizer
	validator *validator.Validate
}

func NewAllTagsHandler(tagRepository TagViewRepository, validate *validator.Validate) AllTagsHandler {
	return AllTagsHandler{tagRepo: tagRepository, sanitizer: newStrictSanitizer(), validator: validate}
}

// Handle performs query to receive all tags for author
func (h AllTagsHandler) Handle(query AllTags) ([]TagView, error) {
	if err := h.validator.Struct(query); err != nil {
		return nil, err
	}

	tags, err := h.tagRepo.GetAllViews(query.AuthorID)

	if err != nil {
		return nil, err
	}

	for i := range tags {
		tags[i].sanitize(h.sanitizer)
	}

	return tags, nil
}
