package translation

import (
	"github.com/google/uuid"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"time"
)

type Translation struct {
	Id            string     `json:"id"`
	AuthorId      string     `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Transcription string     `json:"transcription"`
	Translation   string     `json:"translation"`
	Text          string     `json:"text"`
	Example       string     `json:"example"`
	Tags          []*tag.Tag `json:"tags"`
}

type Request struct {
	Transcription string   `json:"transcription"`
	Translation   string   `json:"translation"`
	Text          string   `json:"text"`
	Example       string   `json:"example"`
	TagIds        []string `json:"tag_ids"`
	AuthorId      string
}

type data struct {
	Request
	Tags []*tag.Tag
}

func newTranslation(data data) *Translation {
	now := time.Now()
	return &Translation{
		Id:            uuid.New().String(),
		AuthorId:      data.AuthorId,
		CreatedAt:     now,
		UpdatedAt:     now,
		Translation:   data.Translation,
		Transcription: data.Transcription,
		Text:          data.Text,
		Example:       data.Example,
		Tags:          data.Tags,
	}
}

func (t *Translation) applyChanges(data data) {
	t.Tags = data.Tags
	t.UpdatedAt = time.Now()
	t.Transcription = data.Transcription
	t.Text = data.Text
	t.Translation = data.Translation
	t.Example = data.Example
}
