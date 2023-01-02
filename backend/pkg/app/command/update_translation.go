package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

// UpdateTranslation update existing translation cmd
type UpdateTranslation struct {
	ID            string
	Transcription string
	Translation   string
	Text          string
	Example       string
	TagIds        []string
	AuthorID      string
}

// UpdateTranslationHandler update existing translation cmd handler
type UpdateTranslationHandler struct {
	translationRepo translation.Repository
	tagRepo         tag.Repository
}

func NewUpdateTranslationHandler(translationRep translation.Repository, tagRepo tag.Repository) UpdateTranslationHandler {
	return UpdateTranslationHandler{
		translationRepo: translationRep,
		tagRepo:         tagRepo,
	}
}

// Handle apply changes from cmd to existing translation
func (h UpdateTranslationHandler) Handle(cmd UpdateTranslation) error {
	tr, err := h.translationRepo.Get(cmd.ID, cmd.AuthorID)

	if err != nil {
		return err
	}

	if err := h.validateTags(cmd); err != nil {
		return err
	}

	tr.ApplyChanges(cmd.Translation, cmd.Transcription, cmd.Text, cmd.Example, cmd.TagIds)

	return h.translationRepo.Update(tr)
}

// validateTags checks that all tags from cmd exist in DB
func (h UpdateTranslationHandler) validateTags(cmd UpdateTranslation) error {
	if len(cmd.TagIds) == 0 {
		return nil
	}

	exists, err := h.tagRepo.AllExist(cmd.TagIds, cmd.AuthorID)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("can not apply changes for translation %s, some passed tag are not found", cmd.ID)
	}

	return nil
}
