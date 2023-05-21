package command //nolint:dupl // it's not fully duplicate

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
)

type AddLang struct {
	Name     string
	AuthorID string
}

type AddLangHandler struct {
	langRepo lang.Repository
}

func NewAddLangHandler(langRepo lang.Repository) AddLangHandler {
	return AddLangHandler{
		langRepo: langRepo,
	}
}

func (h AddLangHandler) Handle(cmd AddLang) (string, error) {
	ln, err := lang.NewLang(cmd.Name, cmd.AuthorID)
	if err != nil {
		return "", err
	}

	if err := h.langRepo.Create(ln); err != nil {
		return "", err
	}

	return ln.ID(), nil
}
