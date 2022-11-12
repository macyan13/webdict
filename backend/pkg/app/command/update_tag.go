package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type UpdateTag struct {
	TagId    string
	Tag      string
	AuthorId string
}

type UpdateTagHandler struct {
	tagRepo tag.Repository
}

func NewUpdateTagHandler(tagRepo tag.Repository) UpdateTagHandler {
	return UpdateTagHandler{tagRepo: tagRepo}
}

func (h UpdateTagHandler) Handle(cmd UpdateTag) error {
	tg := h.tagRepo.GetById(cmd.TagId)

	if !tg.IsAuthor(cmd.AuthorId) {
		return errors.New("can not handle tag update request, translation is not belongs to author")
	}

	tg.ApplyChanges(cmd.Tag)
	return h.tagRepo.Save(*tg)
}
