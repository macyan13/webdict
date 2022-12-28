package query

type LastTranslations struct {
	AuthorID string
	Limit    int
}

type LastTranslationsHandler struct {
	translationRepo TranslationViewRepository
}

func NewLastTranslationsHandler(translationRepo TranslationViewRepository) LastTranslationsHandler {
	return LastTranslationsHandler{translationRepo: translationRepo}
}

func (h LastTranslationsHandler) Handle(cmd LastTranslations) ([]TranslationView, error) {
	limit := cmd.Limit

	if limit == 0 {
		limit = 10
	}

	return h.translationRepo.GetLastViews(cmd.AuthorID, limit)
}
