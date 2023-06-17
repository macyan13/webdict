package query

type RandomTranslations struct {
	AuthorID string
	LangID   string
	TagIds   []string
	Limit    int
}

type RandomTranslationsHandler struct {
	translationRepo TranslationViewRepository
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewRandomTranslationsHandler(translationRepo TranslationViewRepository) RandomTranslationsHandler {
	return RandomTranslationsHandler{translationRepo: translationRepo, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

func (h RandomTranslationsHandler) Handle(query RandomTranslations) (RandomViews, error) {
	limit := h.processLimit(query)
	randomViews, err := h.translationRepo.GetRandomViews(query.AuthorID, query.LangID, query.TagIds, limit)

	if err != nil {
		return randomViews, err
	}

	for i := range randomViews.Views {
		randomViews.Views[i].sanitize(h.strictSntz, h.richSntz)
	}

	return randomViews, nil
}

func (h RandomTranslationsHandler) processLimit(query RandomTranslations) int {
	if query.Limit < 1 || query.Limit > 100 {
		return 10
	}

	return query.Limit
}
