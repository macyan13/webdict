package tag

import (
	"github.com/google/uuid"
)

type Tag struct {
	id       string
	tag      string
	authorID string
}

func NewTag(tag, authorID string) *Tag {
	return &Tag{
		id:       uuid.New().String(),
		tag:      tag,
		authorID: authorID,
	}
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) AuthorID() string {
	return t.authorID
}

func (t *Tag) ApplyChanges(tag string) {
	t.tag = tag
}

func (t *Tag) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       t.id,
		"tag":      t.tag,
		"authorID": t.authorID,
	}
}

func UnmarshalFromDB(
	id string,
	tag string,
	authorID string,
) Tag {
	return Tag{
		id:       id,
		tag:      tag,
		authorID: authorID,
	}
}
