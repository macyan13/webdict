package query

type LastTranslations struct {
	AuthorID string
	Lang     string
	TagIds   []string
	PageSize int
	Page     int
}

type LastTranslationsHandler struct {
	translationRepo TranslationViewRepository
	strictSntz      *strictSanitizer
	richSntz        *richTextSanitizer
}

func NewLastTranslationsHandler(translationRepo TranslationViewRepository) LastTranslationsHandler {
	return LastTranslationsHandler{translationRepo: translationRepo, strictSntz: newStrictSanitizer(), richSntz: newRichTextSanitizer()}
}

func (h LastTranslationsHandler) Handle(query LastTranslations) (LastViews, error) {
	pageSize, page := h.processParameters(query)
	lastViews, err := h.translationRepo.GetLastViews(query.AuthorID, query.Lang, pageSize, page, query.TagIds)

	if err != nil {
		return lastViews, err
	}

	for i := range lastViews.Views {
		lastViews.Views[i].sanitize(h.strictSntz, h.richSntz)
	}

	return lastViews, nil
}

func (h LastTranslationsHandler) processParameters(query LastTranslations) (pageSize, page int) {
	pageSize = query.PageSize

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	page = query.Page

	if page < 1 {
		page = 1
	}

	return
}
