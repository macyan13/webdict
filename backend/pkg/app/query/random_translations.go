package query

import "github.com/go-playground/validator/v10"

type RandomTranslations struct {
	AuthorID string `validate:"required"`
	LangID   string `validate:"required"`
	TagIds   []string
	Limit    int `validate:"gte=1,lte=200"`
}

type RandomTranslationsHandler struct {
	translationRepo TranslationViewRepository
	validator       *validator.Validate
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewRandomTranslationsHandler(translationRepo TranslationViewRepository, validate *validator.Validate) RandomTranslationsHandler {
	return RandomTranslationsHandler{translationRepo: translationRepo, validator: validate, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

func (h RandomTranslationsHandler) Handle(query RandomTranslations) (RandomViews, error) {
	if err := h.validator.Struct(query); err != nil {
		return RandomViews{}, err
	}

	randomViews, err := h.translationRepo.GetRandomViews(query.AuthorID, query.LangID, query.TagIds, query.Limit)

	if err != nil {
		return randomViews, err
	}

	for i := range randomViews.Views {
		randomViews.Views[i].sanitize(h.strictSntz, h.richSntz)
	}

	return randomViews, nil
}
