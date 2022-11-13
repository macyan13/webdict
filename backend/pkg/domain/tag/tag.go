package tag

import (
	"github.com/google/uuid"
	"time"
)

// todo: clean up unused getter after read DB repository implementation
type Tag struct {
	id        string
	createdAt time.Time
	updatedAt time.Time
	tag       string
	authorId  string
}

func NewTag(tag, AuthorId string) *Tag {
	now := time.Now()
	return &Tag{
		id:        uuid.New().String(),
		createdAt: now,
		updatedAt: now,
		tag:       tag,
		authorId:  AuthorId,
	}
}

func (t *Tag) Id() string {
	return t.id
}

func (t *Tag) Tag() string {
	return t.tag
}

func (t *Tag) AuthorId() string {
	return t.authorId
}

func (t *Tag) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Tag) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Tag) ApplyChanges(tag string) {
	t.updatedAt = time.Now()
	t.tag = tag
}

func (t *Tag) IsAuthor(authorId string) bool {
	return t.authorId == authorId
}
