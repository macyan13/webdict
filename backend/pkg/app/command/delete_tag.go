package command

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type DeleteTag struct {
	ID       string
	AuthorID string
}

type DeleteTagHandler struct {
	tagRepo tag.Repository
}

func NewDeleteTagHandler(tagRepo tag.Repository) DeleteTagHandler {
	return DeleteTagHandler{tagRepo: tagRepo}
}

func (h DeleteTagHandler) Handle(cmd DeleteTag) error {
	return h.tagRepo.Delete(cmd.ID, cmd.AuthorID)
}
