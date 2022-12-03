package command

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type AddTag struct {
	Tag      string
	AuthorId string
}

type AddTagHandler struct {
	tagRepo tag.Repository
}

func NewAddTagHandler(tagRepo tag.Repository) AddTagHandler {
	return AddTagHandler{
		tagRepo: tagRepo,
	}
}

func (h AddTagHandler) Handle(cmd AddTag) error {
	return h.tagRepo.Save(*tag.NewTag(cmd.Tag, cmd.AuthorId))
}