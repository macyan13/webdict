package command

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type DeleteTranslation struct {
	ID       string
	AuthorID string
}

type DeleteTranslationHandler struct {
	translationRepo translation.Repository
	tagRepo         tag.Repository
}

func NewDeleteTranslationHandler(translationRepo translation.Repository, tagRepo tag.Repository) DeleteTranslationHandler {
	return DeleteTranslationHandler{
		translationRepo: translationRepo,
		tagRepo:         tagRepo,
	}
}

func (h DeleteTranslationHandler) Handle(cmd DeleteTranslation) error {
	return h.translationRepo.Delete(cmd.ID, cmd.AuthorID)
}
