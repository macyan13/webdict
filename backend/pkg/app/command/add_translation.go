package command

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
)

// AddTranslation create new translation cmd
type AddTranslation struct {
	Transcription string
	Target        string
	Source        string
	Example       string
	TagIds        []string
	AuthorID      string
	LangID        string
}

// AddTranslationHandler create new translation cmd handler
type AddTranslationHandler struct {
	translationRepo translation.Repository
	validator       validator
}

func NewAddTranslationHandler(translationRep translation.Repository, tagRepo tag.Repository, langRepo lang.Repository) AddTranslationHandler {
	return AddTranslationHandler{
		translationRepo: translationRep,
		validator:       newValidator(tagRepo, langRepo),
	}
}

// Handle performs translation creation cmd
func (h AddTranslationHandler) Handle(cmd AddTranslation) (string, error) {
	if err := h.validator.validate(translationData{
		TagIds:   cmd.TagIds,
		LangID:   cmd.LangID,
		AuthorID: cmd.AuthorID,
		Source:   cmd.Source,
	}); err != nil {
		return "", err
	}

	tr, err := translation.NewTranslation(
		cmd.Source,
		cmd.Transcription,
		cmd.Target,
		cmd.AuthorID,
		cmd.Example,
		cmd.TagIds,
		cmd.LangID,
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
