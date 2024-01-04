package query

import (
	"github.com/go-playground/validator/v10"
)

type SearchTranslations struct {
	AuthorID   string `validate:"required"`
	LangID     string `validate:"required"`
	TagIds     []string
	SourcePart string `validate:"excluded_with=TargetPart TagIds"`
	TargetPart string `validate:"excluded_with=SourcePart TagIds"`
	PageSize   int    `validate:"gte=1,lte=200"`
	Page       int    `validate:"gte=1"`
}

type SearchTranslationsHandler struct {
	translationRepo TranslationViewRepository
	validator       *validator.Validate
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewSearchTranslationsHandler(translationRepo TranslationViewRepository, validate *validator.Validate) SearchTranslationsHandler {
	return SearchTranslationsHandler{translationRepo: translationRepo, validator: validate, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

func (h SearchTranslationsHandler) Handle(query SearchTranslations) (LastTranslationViews, error) {
	if err := h.validator.Struct(query); err != nil {
		return LastTranslationViews{}, err
	}

	var lastViews LastTranslationViews
	var err error

	switch {
	case query.SourcePart != "":
		lastViews, err = h.translationRepo.GetLastViewsBySourcePart(query.AuthorID, query.LangID, query.SourcePart, query.PageSize, query.Page)
	case query.TargetPart != "":
		lastViews, err = h.translationRepo.GetLastViewsByTargetPart(query.AuthorID, query.LangID, query.TargetPart, query.PageSize, query.Page)
	default:
		lastViews, err = h.translationRepo.GetLastViewsByTags(query.AuthorID, query.LangID, query.PageSize, query.Page, query.TagIds)
	}

	if err != nil {
		return lastViews, err
	}

	for i := range lastViews.Views {
		lastViews.Views[i].sanitize(h.strictSntz, h.richSntz)
	}

	return lastViews, nil
}
