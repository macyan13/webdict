package tag

import (
	"fmt"
	"github.com/google/uuid"
	"unicode/utf8"
)

type Tag struct {
	id       string
	name     string
	authorID string
}

func NewTag(name, authorID string) (*Tag, error) {
	tg := Tag{
		id:       uuid.New().String(),
		name:     name,
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
	updated := *t
	updated.name = tag

	if err := updated.validate(); err != nil {
		return err
	}

	t.name = tag
	return nil
}

func (t *Tag) validate() error {
	tagCount := utf8.RuneCountInString(t.name)

	if tagCount < 2 {
		return fmt.Errorf("name length should be at least 2 symbols, %d passed (%s)", tagCount, t.name)
	}

	if tagCount > 30 {
		return fmt.Errorf("name max length is 30 symbols, %d passed (%s)", tagCount, t.name)
	}

	if t.authorID == "" {
		return fmt.Errorf("authorID can not be empty")
	}

	return nil
}

func (t *Tag) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       t.id,
		"name":     t.name,
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
		name:     tag,
		authorID: authorID,
	}
}
