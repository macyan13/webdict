package query

type AllLangs struct {
	AuthorID string
}

type AllLangsHandler struct {
	langRepo  LangViewRepository
	sanitizer *strictSanitizer
}

func NewAllLangsHandler(langRepo LangViewRepository) AllLangsHandler {
	return AllLangsHandler{langRepo: langRepo, sanitizer: newStrictSanitizer()}
}

func (h AllLangsHandler) Handle(query AllLangs) ([]LangView, error) {
	langs, err := h.langRepo.GetAllViews(query.AuthorID)

	if err != nil {
		return nil, err
	}

	for i := range langs {
		langs[i].sanitize(h.sanitizer)
	}

	return langs, nil
}
