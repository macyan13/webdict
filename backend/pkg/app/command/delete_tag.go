package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

// DeleteTag delete tag cmd
type DeleteTag struct {
	ID       string
	AuthorID string
}

// DeleteTagHandler Delete tag cmd handler
type DeleteTagHandler struct {
	tagRepo         tag.Repository
	translationRepo translation.Repository
}

func NewDeleteTagHandler(tagRepo tag.Repository, translationRepo translation.Repository) DeleteTagHandler {
	return DeleteTagHandler{tagRepo: tagRepo, translationRepo: translationRepo}
}

// Handle performs tag deletion cmd
func (h *DeleteTagHandler) Handle(cmd DeleteTag) error {
	if err := h.validate(cmd); err != nil {
		return err
	}
	return h.tagRepo.Delete(cmd.ID, cmd.AuthorID)
}

// Validate checks that there is not translation tagged by the tag to be deleted
func (h *DeleteTagHandler) validate(cmd DeleteTag) error {
	exist, err := h.translationRepo.ExistByTag(cmd.ID, cmd.AuthorID)

	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("can not remove tag:%s as some translation is tagged by it", cmd.ID)
	}

	return nil
}
