package command

import (
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
	tg, err := h.tagRepo.Get(cmd.TagId, cmd.AuthorId)

	if err != nil {
		return err
	}

	tg.ApplyChanges(cmd.Tag)
	return h.tagRepo.Update(tg)
}
