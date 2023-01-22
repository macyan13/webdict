package query

// SingleTranslation get translation by ID and authorID query
type SingleTranslation struct {
	ID       string
	AuthorID string
}

// SingleTranslationHandler get translation query handler
type SingleTranslationHandler struct {
	translationRepo TranslationViewRepository
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewSingleTranslationHandler(translationRepo TranslationViewRepository) SingleTranslationHandler {
	return SingleTranslationHandler{translationRepo: translationRepo, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

// Handle performs query to get translation by ID and authorID
func (h SingleTranslationHandler) Handle(cmd SingleTranslation) (TranslationView, error) {
	view, err := h.translationRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return TranslationView{}, err
	}

	view.sanitize(h.strictSntz, h.richSntz)
	return view, nil
}
