package command

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type UpdateTag struct {
	TagID    string
	Tag      string
	AuthorID string
}

type UpdateTagHandler struct {
	tagRepo tag.Repository
}

func NewUpdateTagHandler(tagRepo tag.Repository) UpdateTagHandler {
	return UpdateTagHandler{tagRepo: tagRepo}
}

func (h UpdateTagHandler) Handle(cmd UpdateTag) error {
	tg, err := h.tagRepo.Get(cmd.TagID, cmd.AuthorID)

	if err != nil {
		return err
	}

	tg.ApplyChanges(cmd.Tag)
	return h.tagRepo.Update(tg)
}
