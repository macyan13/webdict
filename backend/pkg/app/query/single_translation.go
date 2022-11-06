package query

type SingleTranslation struct {
	Id       string
	AuthorId string
}

type SingleTranslationRepository interface {
	getTranslation(id, authorId string) *Translation
}

type SingleTranslationHandler struct {
	translationRepo SingleTranslationRepository
}

func NewSingleTranslationHandler(translationRepo SingleTranslationRepository) SingleTranslationHandler {
	return SingleTranslationHandler{translationRepo: translationRepo}
}

func (h SingleTranslationHandler) Handle(cmd SingleTranslation) *Translation {
	return h.translationRepo.getTranslation(cmd.Id, cmd.AuthorId)
}
