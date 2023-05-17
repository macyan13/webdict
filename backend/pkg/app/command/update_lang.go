package command

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
)

type UpdateLang struct {
	ID       string
	Name     string
	AuthorID string
}

type UpdateLangHandler struct {
	langRepo lang.Repository
}

func NewUpdateLangHandler(langRepo lang.Repository) UpdateLangHandler {
	return UpdateLangHandler{langRepo: langRepo}
}

func (h UpdateLangHandler) Handle(cmd UpdateLang) error {
	ln, err := h.langRepo.Get(cmd.ID, cmd.AuthorID)

	if err != nil {
		return err
	}

	if err := ln.ApplyChanges(cmd.Name); err != nil {
		return err
	}

	return h.langRepo.Update(ln)
}
