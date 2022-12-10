package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTag_ApplyChanges(t *testing.T) {
	entity := NewTag("", "")
	tag := "testTag"
	entity.ApplyChanges(tag)
	assert.Equal(t, entity.tag, tag)
}

func TestUnmarshalFromDB(t *testing.T) {
	tag := Tag{
		id:       "testId",
		tag:      "testTag",
		authorId: "testAuthor",
	}

	assert.Equal(t, tag, *UnmarshalFromDB(tag.id, tag.tag, tag.authorId))
}
