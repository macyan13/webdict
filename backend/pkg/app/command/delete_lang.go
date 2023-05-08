package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
)

type DeleteLang struct {
	ID       string
	AuthorID string
}

type DeleteLangHandler struct {
	langRepo        lang.Repository
	translationRepo translation.Repository
}

func NewDeleteLangHandler(langRepo lang.Repository, translationRepo translation.Repository) DeleteLangHandler {
	return DeleteLangHandler{langRepo: langRepo, translationRepo: translationRepo}
}

func (h *DeleteLangHandler) Handle(cmd DeleteLang) error {
	if err := h.validate(cmd); err != nil {
		return err
	}
	return h.langRepo.Delete(cmd.ID, cmd.AuthorID)
}

func (h *DeleteLangHandler) validate(cmd DeleteLang) error {
	exist, err := h.translationRepo.ExistByLang(cmd.ID, cmd.AuthorID)

	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("can not remove lang:%s as some translations use it", cmd.ID)
	}

	return nil
}
