package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
)

// UpdateTranslation update existing translation cmd
type UpdateTranslation struct {
	ID            string
	Source        string
	Transcription string
	Target        string
	AuthorID      string
	Example       string
	TagIds        []string
	Lang          translation.Lang
}

// UpdateTranslationHandler update existing translation cmd handler
type UpdateTranslationHandler struct {
	translationRepo    translation.Repository
	tagRepo            tag.Repository
	supportedLanguages []translation.Lang
}

func NewUpdateTranslationHandler(translationRep translation.Repository, tagRepo tag.Repository, supportedLanguages []translation.Lang) UpdateTranslationHandler {
	return UpdateTranslationHandler{
		translationRepo:    translationRep,
		tagRepo:            tagRepo,
		supportedLanguages: supportedLanguages,
	}
}

// Handle apply changes from cmd to existing translation
func (h UpdateTranslationHandler) Handle(cmd UpdateTranslation) error {
	tr, err := h.translationRepo.Get(cmd.ID, cmd.AuthorID)

	if err != nil {
		return err
	}

	if err := h.validateLang(cmd); err != nil {
		return err
	}

	if err := h.validateTags(cmd); err != nil {
		return err
	}

	if err := tr.ApplyChanges(cmd.Source, cmd.Transcription, cmd.Target, cmd.Example, cmd.TagIds, cmd.Lang); err != nil {
		return err
	}

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

// validateLang check that passed lang is supported
func (h UpdateTranslationHandler) validateLang(cmd UpdateTranslation) error {
	for _, lang := range h.supportedLanguages {
		if cmd.Lang == lang {
			return nil
		}
	}
	return fmt.Errorf("passed language %s is not supported", cmd.Lang)
}
