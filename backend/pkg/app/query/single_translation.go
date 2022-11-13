package query

type SingleTranslation struct {
	Id       string
	AuthorId string
}

type SingleTranslationRepository interface {
	GetTranslation(id, authorId string) *Translation
}

type SingleTranslationHandler struct {
	translationRepo SingleTranslationRepository
}

func NewSingleTranslationHandler(translationRepo SingleTranslationRepository) SingleTranslationHandler {
	return SingleTranslationHandler{translationRepo: translationRepo}
}

func (h SingleTranslationHandler) Handle(cmd SingleTranslation) *Translation {
	return h.translationRepo.GetTranslation(cmd.Id, cmd.AuthorId)
}
