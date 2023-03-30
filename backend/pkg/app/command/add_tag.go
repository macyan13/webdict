package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
)

// AddTag create new tag cmd
type AddTag struct {
	Tag      string
	AuthorID string
}

// AddTagHandler create new tag cmd handler
type AddTagHandler struct {
	tagRepo tag.Repository
}

func NewAddTagHandler(tagRepo tag.Repository) AddTagHandler {
	return AddTagHandler{
		tagRepo: tagRepo,
	}
}

// Handle performs tag creation cmd
func (h AddTagHandler) Handle(cmd AddTag) (string, error) {
	tg, err := tag.NewTag(cmd.Tag, cmd.AuthorID)
	if err != nil {
		return "", err
	}

	if err := h.validate(cmd); err != nil {
		return "", err
	}

	if err := h.tagRepo.Create(tg); err != nil {
		return "", err
	}

	return tg.ID(), nil
}

// Validate checks that there is no such already created tag for the author
func (h AddTagHandler) validate(cmd AddTag) error {
	exist, err := h.tagRepo.ExistByTag(cmd.Tag, cmd.AuthorID)

	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("can not create new tag - tag:%s already created", cmd.Tag)
	}

	return nil
}
