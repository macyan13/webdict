package command

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
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
	LangID        string
}

// UpdateTranslationHandler update existing translation cmd handler
type UpdateTranslationHandler struct {
	translationRepo translation.Repository
	validator       validator
}

func NewUpdateTranslationHandler(translationRep translation.Repository, tagRepo tag.Repository, langRepo lang.Repository) UpdateTranslationHandler {
	return UpdateTranslationHandler{
		translationRepo: translationRep,
		validator:       newValidator(tagRepo, langRepo, translationRep),
	}
}

// Handle apply changes from cmd to existing translation
func (h UpdateTranslationHandler) Handle(cmd UpdateTranslation) error {
	if err := h.validator.validate(translationData{
		TagIds:   cmd.TagIds,
		LangID:   cmd.LangID,
		AuthorID: cmd.AuthorID,
		Source:   cmd.Source,
	}); err != nil {
		return err
	}

	tr, err := h.translationRepo.Get(cmd.ID, cmd.AuthorID)

	if err != nil {
		return err
	}

	if err = tr.ApplyChanges(cmd.Source, cmd.Transcription, cmd.Target, cmd.Example, cmd.TagIds, cmd.LangID); err != nil {
		return err
	}

	return h.translationRepo.Update(tr)
}
