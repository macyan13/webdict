package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type DeleteTranslation struct {
	Id       string
	AuthorId string
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
	tr := h.translationRepo.GetById(cmd.Id)

	if !tr.IsAuthor(cmd.AuthorId) {
		return errors.New("can not handle translation delete request, translation is not belongs to author")
	}

	return h.translationRepo.Delete(*tr)
}
