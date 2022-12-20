package tag

import (
	"github.com/google/uuid"
)

type Tag struct {
	id       string
	tag      string
	authorId string
}

func NewTag(tag, AuthorId string) *Tag {
	return &Tag{
		id:       uuid.New().String(),
		tag:      tag,
		authorId: AuthorId,
	}
}

func (t *Tag) Id() string {
	return t.id
}

func (t *Tag) AuthorId() string {
	return t.authorId
}

func (t *Tag) ApplyChanges(tag string) {
	t.tag = tag
}

func (t *Tag) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       t.id,
		"tag":      t.tag,
		"authorId": t.authorId,
	}
}

func UnmarshalFromDB(
	id string,
	tag string,
	authorId string,
) Tag {
	return Tag{
		id:       id,
		tag:      tag,
		authorId: authorId,
	}
}
