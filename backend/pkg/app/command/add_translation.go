package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

// AddTranslation create new translation cmd
type AddTranslation struct {
	Transcription string
	Translation   string
	Text          string
	Example       string
	TagIds        []string
	AuthorID      string
}

// AddTranslationHandler create new translation cmd handler
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

// Handle performs translation creation cmd
func (h AddTranslationHandler) Handle(cmd AddTranslation) (string, error) {
	if err := h.validateTags(cmd); err != nil {
		return "", err
	}

	if err := h.validateTranslation(cmd); err != nil {
		return "", err
	}

	tr, err := translation.NewTranslation(
		cmd.Text,
		cmd.Transcription,
		cmd.Translation,
		cmd.AuthorID,
		cmd.Example,
		cmd.TagIds,
	)
	if err != nil {
		return "", err
	}

	err = h.translationRepo.Create(tr)

	if err != nil {
		return "", err
	}

	return tr.ID(), nil
}

// validateTags check that all cmd tags exist
func (h AddTranslationHandler) validateTags(cmd AddTranslation) error {
	if len(cmd.TagIds) == 0 {
		return nil
	}

	exist, err := h.tagRepo.AllExist(cmd.TagIds, cmd.AuthorID)

	if err != nil {
		return err
	}
	if !exist {
		return errors.New("can not apply changes for translation tags, some passed tag are not found")
	}

	return nil
}

// validateTranslation checks that there is not already creation translation with the cmd text
func (h AddTranslationHandler) validateTranslation(cmd AddTranslation) error {
	exist, err := h.translationRepo.ExistByText(cmd.Text, cmd.AuthorID)

	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("translation with text: %s already created", cmd.Text)
	}

	return nil
}
