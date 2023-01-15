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
		views[i].sanitise()
	}

	return views, nil
}
