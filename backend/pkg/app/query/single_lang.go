package query

import "github.com/go-playground/validator/v10"

type SingleLang struct {
	ID       string `validate:"required"`
	AuthorID string `validate:"required"`
}

type SingleLangHandler struct {
	langRepo   LangViewRepository
	validator  *validator.Validate
	strictSntz *strictSanitizer
}

func NewSingleLangHandler(langRepo LangViewRepository, validate *validator.Validate) SingleLangHandler {
	return SingleLangHandler{langRepo: langRepo, validator: validate, strictSntz: newStrictSanitizer()}
}

func (h SingleLangHandler) Handle(cmd SingleLang) (LangView, error) {
	if err := h.validator.Struct(cmd); err != nil {
		return LangView{}, err
	}
	view, err := h.langRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return LangView{}, err
	}

	view.sanitize(h.strictSntz)
	return view, nil
}
