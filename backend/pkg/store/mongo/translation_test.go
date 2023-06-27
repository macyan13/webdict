package mongo

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTranslationRepo_fromDomainToModel(t *testing.T) {
	meaning := "testTranslation"
	authorID := "testAuthor"
	transcription := "testTranscription"
	source := "testText"
	example := "testExample"
	tags := []string{"test1", "test2"}
	langID := "EN"
	domain, err := translation.NewTranslation(source, transcription, meaning, authorID, example, tags, langID)
	assert.Nil(t, err)
	domainMap := domain.ToMap()

	repo := TranslationRepo{}
	model, err := repo.fromDomainToModel(domain)

	assert.Nil(t, err)
	assert.Equal(t, meaning, model.Target)
	assert.Equal(t, authorID, model.AuthorID)
	assert.Equal(t, transcription, model.Transcription)
	assert.Equal(t, source, model.Source)
	assert.Equal(t, example, model.Example)
	assert.Equal(t, tags, model.TagIDs)
	assert.Equal(t, langID, model.LangID)
	assert.Equal(t, domainMap["createdAt"], model.CreatedAt)
	assert.Equal(t, domainMap["updatedAt"], model.UpdatedAt)
}

func TestTranslationRepo_fromModelToView_positiveCase(t *testing.T) {
	model := TranslationModel{
		ID:            "id",
		AuthorID:      "testAuthor",
		CreatedAt:     time.Now().Add(5 * time.Second),
		UpdatedAt:     time.Now().Add(10 * time.Second),
		Transcription: "transcription",
		Target:        "translation",
		Source:        "text",
		Example:       "example",
		TagIDs:        []string{"tag1", "tag2"},
		LangID:        "EN",
	}

	tagViews := []query.TagView{{Name: "tag1"}, {Name: "tag2"}}
	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return(tagViews, nil)

	langRepo := query.MockLangViewRepository{}
	langRepo.On("GetView", "EN", "testAuthor").Return(query.LangView{ID: "EN", Name: "English"}, nil)

	translationRepo := &TranslationRepo{
		tagRepo:  &tagRepo,
		langRepo: &langRepo,
	}

	view, err := translationRepo.fromModelToView(model)
	assert.Nil(t, err)

	assert.Equal(t, model.ID, view.ID)
	assert.Equal(t, model.CreatedAt, view.CreatedAd)
	assert.Equal(t, model.Target, view.Target)
	assert.Equal(t, model.Transcription, view.Transcription)
	assert.Equal(t, model.Source, view.Source)
	assert.Equal(t, model.Example, view.Example)
	assert.Equal(t, tagViews, view.Tags)
	assert.Equal(t, model.LangID, view.Lang.ID)
}

func TestTranslationRepo_fromModelToView_tagsNotSet(t *testing.T) {
	model := TranslationModel{
		ID:            "id",
		AuthorID:      "testAuthor",
		CreatedAt:     time.Now().Add(5 * time.Second),
		UpdatedAt:     time.Now().Add(10 * time.Second),
		Transcription: "transcription",
		Target:        "translation",
		Source:        "text",
		Example:       "example",
		LangID:        "EN",
	}

	langRepo := query.MockLangViewRepository{}
	langRepo.On("GetView", "EN", "testAuthor").Return(query.LangView{ID: "EN"}, nil)

	translationRepo := &TranslationRepo{langRepo: &langRepo}

	view, err := translationRepo.fromModelToView(model)
	assert.Nil(t, err)
	assert.Equal(t, model.ID, view.ID)
	assert.Equal(t, model.CreatedAt, view.CreatedAd)
	assert.Equal(t, model.Target, view.Target)
	assert.Equal(t, model.Transcription, view.Transcription)
	assert.Equal(t, model.Source, view.Source)
	assert.Equal(t, model.Example, view.Example)
	assert.Equal(t, model.LangID, view.Lang.ID)
}

func TestTranslationRepo_fromModelToView_errorOnGetViews(t *testing.T) {
	model := TranslationModel{
		AuthorID: "testAuthor",
		TagIDs:   []string{"tag1", "tag2"},
		LangID:   "EN",
	}

	langRepo := query.MockLangViewRepository{}
	langRepo.On("GetView", "EN", "testAuthor").Return(query.LangView{ID: "EN"}, nil)

	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return(nil, fmt.Errorf("dbError"))

	translationRepo := &TranslationRepo{
		tagRepo:  &tagRepo,
		langRepo: &langRepo,
	}

	_, err := translationRepo.fromModelToView(model)
	assert.Equal(t, "dbError", err.Error())
}

func TestTranslationRepo_fromModelToView_errorOnTagViewsMiscount(t *testing.T) {
	model := TranslationModel{
		AuthorID: "testAuthor",
		TagIDs:   []string{"tag1", "tag2"},
		LangID:   "EN",
	}

	langRepo := query.MockLangViewRepository{}
	langRepo.On("GetView", "EN", "testAuthor").Return(query.LangView{ID: "EN"}, nil)

	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return([]query.TagView{{Name: "tag1"}}, nil)

	translationRepo := &TranslationRepo{
		tagRepo:  &tagRepo,
		langRepo: &langRepo,
	}

	_, err := translationRepo.fromModelToView(model)
	assert.Equal(t, "can not find all translation tags", err.Error())
}

func TestTranslationRepo_fromModelToView_errorOnGetLangView(t *testing.T) {
	model := TranslationModel{
		AuthorID: "testAuthor",
		TagIDs:   []string{"tag1"},
		LangID:   "EN",
	}

	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1"}, "testAuthor").Return([]query.TagView{{Name: "tag1"}}, nil)

	langRepo := query.MockLangViewRepository{}
	langRepo.On("GetView", "EN", "testAuthor").Return(query.LangView{}, errors.New("testError"))

	translationRepo := &TranslationRepo{
		tagRepo:  &tagRepo,
		langRepo: &langRepo,
	}

	_, err := translationRepo.fromModelToView(model)
	assert.Equal(t, "testError", err.Error())
}
