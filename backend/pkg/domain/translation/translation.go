package translation

import (
	"github.com/google/uuid"
	"time"
)

type Translation struct {
	id            string
	authorId      string
	createdAt     time.Time
	updatedAt     time.Time
	transcription string
	translation   string
	text          string
	example       string
	tagIds        []string
}

func NewTranslation(translation, transcription, text, example, authorId string, tagIds []string) *Translation {
	now := time.Now()
	return &Translation{
		id:            uuid.New().String(),
		authorId:      authorId,
		createdAt:     now,
		updatedAt:     now,
		translation:   translation,
		transcription: transcription,
		text:          text,
		example:       example,
		tagIds:        tagIds,
	}
}

func (t *Translation) Id() string {
	return t.id
}

func (t *Translation) AuthorId() string {
	return t.authorId
}

func (t *Translation) ApplyChanges(translation, transcription, text, example string, tagIds []string) {
	t.tagIds = tagIds
	t.transcription = transcription
	t.text = text
	t.translation = translation
	t.example = example
	t.updatedAt = time.Now()
}

func (t *Translation) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            t.id,
		"authorId":      t.authorId,
		"createdAt":     t.createdAt,
		"updatedAt":     t.updatedAt,
		"translation":   t.translation,
		"transcription": t.transcription,
		"text":          t.text,
		"example":       t.example,
		"tagIds":        t.tagIds,
	}
}

func UnmarshalFromDB(
	id string,
	authorId string,
	createdAt time.Time,
	updatedAt time.Time,
	transcription string,
	translation string,
	text string,
	example string,
	tagIds []string,
) *Translation {
	return &Translation{
		id:            id,
		authorId:      authorId,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		transcription: transcription,
		translation:   translation,
		text:          text,
		example:       example,
		tagIds:        tagIds,
	}
}
