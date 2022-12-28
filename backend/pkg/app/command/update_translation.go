package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type UpdateTranslation struct {
	ID            string
	Transcription string
	Translation   string
	Text          string
	Example       string
	TagIds        []string
	AuthorID      string
}

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

func (h UpdateTranslationHandler) Handle(cmd *UpdateTranslation) error {
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

func (h UpdateTranslationHandler) validateTags(cmd *UpdateTranslation) error {
	if len(cmd.TagIds) == 0 {
		return nil
	}

	exists, err := h.tagRepo.AllExist(cmd.TagIds, cmd.AuthorID)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("can not apply changes for translation tags, some passed tag are not found")
	}

	return nil
}
