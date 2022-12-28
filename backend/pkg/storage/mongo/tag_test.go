package mongo

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagRepo_fromDomainToModel(t *testing.T) {
	tagValue := "testTag"
	entity := tag.NewTag(tagValue, "testAuthor")
	repo := TagRepo{}

	model, err := repo.fromDomainToModel(*entity)
	assert.Nil(t, err)
	assert.Equal(t, entity.ID(), model.ID)
	assert.Equal(t, entity.AuthorID(), model.AuthorID)
	assert.Equal(t, tagValue, model.Tag)
}

func TestTagRepo_fromModelToView(t *testing.T) {
	model := TagModel{
		ID:       "id",
		Tag:      "tag",
		AuthorID: "author",
	}

	repo := TagRepo{}
	view := repo.fromModelToView(model)
	assert.Equal(t, model.ID, view.ID)
	assert.Equal(t, model.Tag, view.Tag)
}
