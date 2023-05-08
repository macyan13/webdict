package query

type SingleLang struct {
	ID       string
	AuthorID string
}

type SingleLangHandler struct {
	langRepo   LangViewRepository
	strictSntz *strictSanitizer
}

func NewSingleLangHandler(langRepo LangViewRepository) SingleLangHandler {
	return SingleLangHandler{langRepo: langRepo, strictSntz: newStrictSanitizer()}
}

func (h SingleLangHandler) Handle(cmd SingleLang) (LangView, error) {
	view, err := h.langRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return LangView{}, err
	}

	view.sanitize(h.strictSntz)
	return view, nil
}
