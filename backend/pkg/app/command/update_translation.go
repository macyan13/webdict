package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type UpdateTranslation struct {
	Id            string
	Transcription string
	Translation   string
	Text          string
	Example       string
	TagIds        []string
	AuthorId      string
}

type UpdateTranslationHandler struct {
	translationRep translation.Repository
	tagRepo        tag.Repository
}

func (h UpdateTranslationHandler) Handle(cmd UpdateTranslation) error {
	tr := h.translationRep.GetById(cmd.Id)

	if !tr.IsAuthor(cmd.AuthorId) {
		return errors.New("can not handle translation update request, translation is not belongs to author")
	}

	if err := h.validateTags(cmd); err != nil {
		return err
	}

	tr.ApplyChanges(cmd.Translation, cmd.Transcription, cmd.Text, cmd.Example, cmd.TagIds)

	return h.translationRep.Save(*tr)
}

func (h UpdateTranslationHandler) validateTags(cmd UpdateTranslation) error {
	if len(cmd.TagIds) == 0 {
		return nil
	}

	if !h.tagRepo.AllExist(cmd.TagIds, cmd.AuthorId) {
		return errors.New("can not apply changes for translation tags, some passed tag are not found")
	}

	return nil
}
