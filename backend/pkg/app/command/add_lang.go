package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
)

var ErrLangAlreadyExists = errors.New("lang already exists")

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

	if err := h.validate(cmd); err != nil {
		return "", err
	}

	if err := h.langRepo.Create(ln); err != nil {
		return "", err
	}

	return ln.ID(), nil
}

func (h AddLangHandler) validate(cmd AddLang) error {
	exist, err := h.langRepo.ExistByName(cmd.Name, cmd.AuthorID)

	if err != nil {
		return err
	}

	if exist {
		return ErrLangAlreadyExists
	}

	return nil
}
