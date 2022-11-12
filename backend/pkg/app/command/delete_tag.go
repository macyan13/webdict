package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type DeleteTag struct {
	Id       string
	AuthorId string
}

type DeleteTagHandler struct {
	tagRepo tag.Repository
}

func NewDeleteTagHandler(tagRepo tag.Repository) DeleteTagHandler {
	return DeleteTagHandler{tagRepo: tagRepo}
}

func (h DeleteTagHandler) Handle(cmd DeleteTag) error {
	tg := h.tagRepo.GetById(cmd.Id)

	if !tg.IsAuthor(cmd.AuthorId) {
		return errors.New("can not handle tag delete request, translation is not belongs to author")
	}

	return h.tagRepo.Delete(*tg)
}
