package tag

import (
	"fmt"
	"github.com/google/uuid"
)

type Tag struct {
	id       string
	tag      string
	authorID string
}

func NewTag(tag, authorID string) (*Tag, error) {
	tg := Tag{
		id:       uuid.New().String(),
		tag:      tag,
		authorID: authorID,
	}

	if err := tg.validate(); err != nil {
		return nil, err
	}

	return &tg, nil
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) AuthorID() string {
	return t.authorID
}

func (t *Tag) ApplyChanges(tag string) error {
	t.tag = tag
	return t.validate()
}

func (t *Tag) validate() error {
	if len(t.tag) < 2 {
		return fmt.Errorf("tag length should be at least 2 symbols, %d passed", len(t.tag))
	}

	if len(t.tag) > 30 {
		return fmt.Errorf("tag max length is 30 symbols, %d passed", len(t.tag))
	}

	if t.authorID == "" {
		return fmt.Errorf("authorID can not be empty")
	}

	return nil
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
) *Tag {
	return &Tag{
		id:       id,
		tag:      tag,
		authorID: authorID,
	}
}
