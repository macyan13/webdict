package mongo

import (
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLangRepo_fromModelToView(t *testing.T) {
	model := LangModel{
		ID:       "id",
		Name:     "en",
		AuthorID: "author",
	}

	repo := LangRepo{}
	view := repo.fromModelToView(model)
	assert.Equal(t, model.ID, view.ID)
	assert.Equal(t, model.Name, view.Name)
}

func TestLangRepo_fromDomainToModel(t *testing.T) {
	ln := "en"
	entity, err := lang.NewLang(ln, "testAuthor")
	assert.Nil(t, err)
	repo := LangRepo{}

	model, err := repo.fromDomainToModel(entity)
	assert.Nil(t, err)
	assert.Equal(t, entity.ID(), model.ID)
	assert.Equal(t, entity.AuthorID(), model.AuthorID)
	assert.Equal(t, ln, model.Name)
}
