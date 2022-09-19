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
}

type Request struct {
	Tag string `json:"tag"`
}

func NewTag(request Request) *Tag {
	now := time.Now().Unix()
	return &Tag{
		Id:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Tag:       request.Tag,
	}
}

func (t *Tag) ApplyChanges(request Request) {
	t.UpdatedAt = time.Now().Unix()
	t.Tag = request.Tag
}
