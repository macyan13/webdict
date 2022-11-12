package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTag(t *testing.T) {
	tag := NewTag(Request{})
	assert.Equal(t, tag.UpdatedAt, tag.CreatedAt)
}

func TestTag_ApplyChanges(t *testing.T) {
	entity := NewTag(Request{})
	updatedAt := entity.UpdatedAt

	tag := "tag"
	request := Request{
		Tag: tag,
	}

	time.Sleep(time.Second)

	entity.ApplyChanges(request)

	assert := assert.New(t)
	assert.Equal(entity.tag, tag)
	assert.Greaterf(entity.UpdatedAt, updatedAt, "error message %s", "formatted")
}
