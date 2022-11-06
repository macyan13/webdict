package translation

import (
	"github.com/google/uuid"
	"time"
)

type Translation struct {
	Id            string    `json:"id"`
	AuthorId      string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Transcription string    `json:"transcription"`
	Translation   string    `json:"translation"`
	Text          string    `json:"text"`
	Example       string    `json:"example"`
	TagIds        []string  `json:"tags"`
}

type Request struct {
	Transcription string   `json:"transcription"`
	Translation   string   `json:"translation"`
	Text          string   `json:"text"`
	Example       string   `json:"example"`
	TagIds        []string `json:"tag_ids"`
	AuthorId      string
}

func NewTranslation(translation, transcription, text, example, authorId string, tagIds []string) *Translation {
	now := time.Now()
	return &Translation{
		Id:            uuid.New().String(),
		AuthorId:      authorId,
		CreatedAt:     now,
		UpdatedAt:     now,
		Translation:   translation,
		Transcription: transcription,
		Text:          text,
		Example:       example,
		TagIds:        tagIds,
	}
}

func (t *Translation) ApplyChanges(translation, transcription, text, example string, tagIds []string) {
	t.TagIds = tagIds
	t.Transcription = transcription
	t.Text = text
	t.Translation = translation
	t.Example = example
	t.UpdatedAt = time.Now()
}

func (t *Translation) IsAuthor(authorId string) bool {
	return t.AuthorId == authorId
}
