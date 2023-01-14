package mongo

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTranslationRepo_fromDomainToModel(t *testing.T) {
	tr := "testTranslation"
	authorID := "testAuthor"
	transcription := "testTranscription"
	text := "testText"
	example := "testExample"
	tags := []string{"test1", "test2"}
	domain := translation.NewTranslation(text, transcription, tr, authorID, example, tags)
	domainMap := domain.ToMap()

	repo := TranslationRepo{}
	model, err := repo.fromDomainToModel(domain)

	assert.Nil(t, err)
	assert.Equal(t, tr, model.Translation)
	assert.Equal(t, authorID, model.AuthorID)
	assert.Equal(t, transcription, model.Transcription)
	assert.Equal(t, text, model.Text)
	assert.Equal(t, example, model.Example)
	assert.Equal(t, tags, model.TagIDs)
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
		Translation:   "translation",
		Text:          "text",
		Example:       "example",
		TagIDs:        []string{"tag1", "tag2"},
	}

	tagViews := []query.TagView{{Tag: "tag1"}, {Tag: "tag2"}}
	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return(tagViews, nil)

	translationRepo := &TranslationRepo{
		tagRepo: &tagRepo,
	}

	view, err := translationRepo.fromModelToView(model)
	assert.Nil(t, err)

	assert.Equal(t, model.ID, view.ID)
	assert.Equal(t, model.CreatedAt, view.CreatedAd)
	assert.Equal(t, model.Translation, view.Translation)
	assert.Equal(t, model.Transcription, view.Transcription)
	assert.Equal(t, model.Text, view.Text)
	assert.Equal(t, model.Example, view.Example)
	assert.Equal(t, tagViews, view.Tags)
}

func TestTranslationRepo_fromModelToView_errorOnGetViews(t *testing.T) {
	model := TranslationModel{
		AuthorID: "testAuthor",
		TagIDs:   []string{"tag1", "tag2"},
	}

	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return(nil, fmt.Errorf("dbError"))

	translationRepo := &TranslationRepo{
		tagRepo: &tagRepo,
	}

	_, err := translationRepo.fromModelToView(model)
	assert.Equal(t, "dbError", err.Error())
}

func TestTranslationRepo_fromModelToView_errorOnTagViewsMiscount(t *testing.T) {
	model := TranslationModel{
		AuthorID: "testAuthor",
		TagIDs:   []string{"tag1", "tag2"},
	}

	tagRepo := query.MockTagViewRepository{}
	tagRepo.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return([]query.TagView{{Tag: "tag1"}}, nil)

	translationRepo := &TranslationRepo{
		tagRepo: &tagRepo,
	}

	_, err := translationRepo.fromModelToView(model)
	assert.Equal(t, "can not find all translation tags", err.Error())
}
