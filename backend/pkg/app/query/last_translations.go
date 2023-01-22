package query

type LastTranslations struct {
	AuthorID string
	Limit    int
}

type LastTranslationsHandler struct {
	translationRepo TranslationViewRepository
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewLastTranslationsHandler(translationRepo TranslationViewRepository) LastTranslationsHandler {
	return LastTranslationsHandler{translationRepo: translationRepo, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

// Handle todo: add test after finalizing solution
func (h LastTranslationsHandler) Handle(cmd LastTranslations) ([]TranslationView, error) {
	limit := cmd.Limit

	if limit == 0 {
		limit = 10
	}

	views, err := h.translationRepo.GetLastViews(cmd.AuthorID, limit)

	if err != nil {
		return nil, err
	}

	for i := range views {
		views[i].sanitize(h.strictSntz, h.richSntz)
	}

	return views, nil
}
