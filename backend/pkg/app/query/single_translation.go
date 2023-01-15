package query

// SingleTranslation get translation by ID and authorID query
type SingleTranslation struct {
	ID       string
	AuthorID string
}

// SingleTranslationHandler get translation query handler
type SingleTranslationHandler struct {
	translationRepo TranslationViewRepository
}

func NewSingleTranslationHandler(translationRepo TranslationViewRepository) SingleTranslationHandler {
	return SingleTranslationHandler{translationRepo: translationRepo}
}

// Handle performs query to get translation by ID and authorID
func (h SingleTranslationHandler) Handle(cmd SingleTranslation) (TranslationView, error) {
	view, err := h.translationRepo.GetView(cmd.ID, cmd.AuthorID)

	if err != nil {
		return TranslationView{}, err
	}

	view.sanitise()
	return view, nil
}
