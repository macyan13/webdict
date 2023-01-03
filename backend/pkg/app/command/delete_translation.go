package command

import (
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

// DeleteTranslation cmd
type DeleteTranslation struct {
	ID       string
	AuthorID string
}

// DeleteTranslationHandler delete translation cmd handler
type DeleteTranslationHandler struct {
	translationRepo translation.Repository
}

func NewDeleteTranslationHandler(translationRepo translation.Repository) DeleteTranslationHandler {
	return DeleteTranslationHandler{
		translationRepo: translationRepo,
	}
}

// Handle performs translation deletion
func (h DeleteTranslationHandler) Handle(cmd DeleteTranslation) error {
	return h.translationRepo.Delete(cmd.ID, cmd.AuthorID)
}
