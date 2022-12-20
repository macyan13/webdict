package command

import (
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
	return h.tagRepo.Delete(cmd.Id, cmd.AuthorId)
}
