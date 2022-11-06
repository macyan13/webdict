package tag

import (
	"github.com/google/uuid"
	"time"
)

type Tag struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Tag       string `json:"tag"`
	AuthorId  string
}

type Request struct {
	Tag string `json:"tag"`
}

func NewTag(tag, AuthorId string) *Tag {
	now := time.Now().Unix()
	return &Tag{
		Id:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Tag:       tag,
		AuthorId:  AuthorId,
	}
}

func (t *Tag) ApplyChanges(tag string) {
	t.UpdatedAt = time.Now().Unix()
	t.Tag = tag
}

func (t *Tag) IsAuthor(authorId string) bool {
	return t.AuthorId == authorId
}
