package query

import "github.com/go-playground/validator/v10"

// SingleTranslation get translation by ID and authorID query
type SingleTranslation struct {
	ID       string `validate:"required"`
	AuthorID string `validate:"required"`
}

// SingleTranslationHandler get translation query handler
type SingleTranslationHandler struct {
	translationRepo TranslationViewRepository
	validator       *validator.Validate
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewSingleTranslationHandler(translationRepo TranslationViewRepository, validate *validator.Validate) SingleTranslationHandler {
	return SingleTranslationHandler{translationRepo: translationRepo, validator: validate, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

// Handle performs query to get translation by ID and authorID
func (h SingleTranslationHandler) Handle(cmd SingleTranslation) (TranslationView, error) {
	if err := h.validator.Struct(cmd); err != nil {
		return TranslationView{}, err
	}

	view, err := h.translationRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return TranslationView{}, err
	}

	view.sanitize(h.strictSntz, h.richSntz)
	return view, nil
}
