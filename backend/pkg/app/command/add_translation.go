package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type AddTranslation struct {
	Transcription string   `json:"transcription"`
	Translation   string   `json:"translation"`
	Text          string   `json:"text"`
	Example       string   `json:"example"`
	TagIds        []string `json:"tag_ids"`
	AuthorId      string
}

type AddTranslationHandler struct {
	translationRepo translation.Repository
	tagRepo         tag.Repository
}

func NewAddTranslationHandler(translationRep translation.Repository, tagRepo tag.Repository) AddTranslationHandler {
	return AddTranslationHandler{
		translationRepo: translationRep,
		tagRepo:         tagRepo,
	}
}

func (h AddTranslationHandler) Handle(cmd AddTranslation) error {
	if err := h.validateTags(cmd); err != nil {
		return err
	}

	tr := translation.NewTranslation(
		cmd.Translation,
		cmd.Transcription,
		cmd.Text,
		cmd.Example,
		cmd.AuthorId,
		cmd.TagIds,
	)

	return h.translationRepo.Save(*tr)
}

func (h AddTranslationHandler) validateTags(cmd AddTranslation) error {
	if len(cmd.TagIds) == 0 {
		return nil
	}

	if !h.tagRepo.AllExist(cmd.TagIds, cmd.AuthorId) {
		return errors.New("can not apply changes for translation tags, some passed tag are not found")
	}

	return nil
}