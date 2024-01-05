package query

import "github.com/go-playground/validator/v10"

type AllLangs struct {
	AuthorID string `validate:"required"`
}

type AllLangsHandler struct {
	langRepo  LangViewRepository
	sanitizer *strictSanitizer
	validator *validator.Validate
}

func NewAllLangsHandler(langRepo LangViewRepository, validate *validator.Validate) AllLangsHandler {
	return AllLangsHandler{langRepo: langRepo, sanitizer: newStrictSanitizer(), validator: validate}
}

func (h AllLangsHandler) Handle(query AllLangs) ([]LangView, error) {
	if err := h.validator.Struct(query); err != nil {
		return nil, err
	}

	langs, err := h.langRepo.GetAllViews(query.AuthorID)

	if err != nil {
		return nil, err
	}

	for i := range langs {
		langs[i].sanitize(h.sanitizer)
	}

	return langs, nil
}
