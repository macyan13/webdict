package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
)

// AddTranslation create new translation cmd
type AddTranslation struct {
	Transcription string
	Target        string
	Text          string
	Example       string
	TagIds        []string
	AuthorID      string
	Lang          translation.Lang
}

// AddTranslationHandler create new translation cmd handler
type AddTranslationHandler struct {
	translationRepo    translation.Repository
	tagRepo            tag.Repository
	supportedLanguages []translation.Lang
}

func NewAddTranslationHandler(
	translationRep translation.Repository,
	tagRepo tag.Repository,
	supportedLanguages []translation.Lang,
) AddTranslationHandler {
	return AddTranslationHandler{
		translationRepo:    translationRep,
		tagRepo:            tagRepo,
		supportedLanguages: supportedLanguages,
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

	if err := h.validateLang(cmd); err != nil {
		return "", err
	}

	tr, err := translation.NewTranslation(
		cmd.Text,
		cmd.Transcription,
		cmd.Target,
		cmd.AuthorID,
		cmd.Example,
		cmd.TagIds,
		cmd.Lang,
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
		return errors.New("some passed tag are not found")
	}

	return nil
}

// validateLang check that passed lang is supported
func (h AddTranslationHandler) validateLang(cmd AddTranslation) error {
	for _, lang := range h.supportedLanguages {
		if cmd.Lang == lang {
			return nil
		}
	}
	return fmt.Errorf("passed language %s is not supported", cmd.Lang)
}

// validateTranslation checks that there is not already creation translation with the cmd text
func (h AddTranslationHandler) validateTranslation(cmd AddTranslation) error {
	exist, err := h.translationRepo.ExistBySource(cmd.Text, cmd.AuthorID)

	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("translation with text: %s already created", cmd.Text)
	}

	return nil
}
