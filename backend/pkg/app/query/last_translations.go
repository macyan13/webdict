package query

type LastTranslations struct {
	AuthorId string
	Limit    int
}

type LastTranslationsRepository interface {
	getLastTranslations(authorId string, limit int) []Translation
}

type LastTranslationsHandler struct {
	translationRepo LastTranslationsRepository
}

func NewLastTranslationsHandler(translationRepo LastTranslationsRepository) LastTranslationsHandler {
	return LastTranslationsHandler{translationRepo: translationRepo}
}

func (h LastTranslationsHandler) Handle(cmd LastTranslations) []Translation {
	limit := cmd.Limit

	if limit == 0 {
		limit = 10
	}

	return h.translationRepo.getLastTranslations(cmd.AuthorId, limit)
}
