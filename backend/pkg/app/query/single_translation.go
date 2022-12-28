package query

type SingleTranslation struct {
	ID       string
	AuthorID string
}

type SingleTranslationHandler struct {
	translationRepo TranslationViewRepository
}

func NewSingleTranslationHandler(translationRepo TranslationViewRepository) SingleTranslationHandler {
	return SingleTranslationHandler{translationRepo: translationRepo}
}

func (h SingleTranslationHandler) Handle(cmd SingleTranslation) (TranslationView, error) {
	return h.translationRepo.GetView(cmd.ID, cmd.AuthorID)
}
