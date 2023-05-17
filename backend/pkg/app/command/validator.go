package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
)

type validator struct {
	tagRepo         tag.Repository
	langRepo        lang.Repository
	translationRepo translation.Repository
}

type translationData struct {
	TagIds   []string
	LangID   string
	AuthorID string
	Source   string
}

func newValidator(tagRepo tag.Repository, langRepo lang.Repository, translationRepo translation.Repository) validator {
	return validator{
		tagRepo:         tagRepo,
		langRepo:        langRepo,
		translationRepo: translationRepo,
	}
}

func (v validator) validate(data translationData) error {
	var err error

	err = errors.Join(err, v.validateTags(data))
	err = errors.Join(err, v.validateLang(data))
	err = errors.Join(err, v.validateSource(data))

	return err
}

func (v validator) validateTags(data translationData) error {
	if len(data.TagIds) == 0 {
		return nil
	}

	exist, err := v.tagRepo.AllExist(data.TagIds, data.AuthorID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("some of passed tags: %v are not found", data.TagIds)
	}

	return nil
}

func (v validator) validateLang(data translationData) error {
	exist, err := v.langRepo.Exist(data.LangID, data.AuthorID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("lang with id: %s is not found", data.LangID)
	}

	return nil
}

func (v validator) validateSource(data translationData) error {
	exist, err := v.translationRepo.ExistBySource(data.Source, data.AuthorID, data.LangID)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("translation with source: %s already created", data.Source)
	}

	return nil
}
