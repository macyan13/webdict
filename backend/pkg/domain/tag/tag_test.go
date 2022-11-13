package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTag(t *testing.T) {
	tag := NewTag("", "")
	assert.Equal(t, tag.updatedAt, tag.createdAt)
}

func TestTag_ApplyChanges(t *testing.T) {
	entity := NewTag("", "")
	updatedAt := entity.updatedAt
	tag := "testTag"

	time.Sleep(time.Second)
	entity.ApplyChanges(tag)

	assert := assert.New(t)
	assert.Equal(entity.tag, tag)
	assert.Greaterf(entity.updatedAt, updatedAt, "Tag.ApplyChanges - updatedAt should be greater createdAt")
}
