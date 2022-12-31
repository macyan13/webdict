package translation

import (
	"github.com/google/uuid"
	"time"
)

type Translation struct {
	id            string
	authorID      string
	createdAt     time.Time
	updatedAt     time.Time
	transcription string
	translation   string
	text          string
	example       string
	tagIDs        []string
}

func NewTranslation(translation, transcription, text, example, authorID string, tagIDs []string) *Translation {
	now := time.Now()
	return &Translation{
		id:            uuid.New().String(),
		authorID:      authorID,
		createdAt:     now,
		updatedAt:     now,
		translation:   translation,
		transcription: transcription,
		text:          text,
		example:       example,
		tagIDs:        tagIDs,
	}
}

func (t *Translation) ID() string {
	return t.id
}

func (t *Translation) AuthorID() string {
	return t.authorID
}

func (t *Translation) ApplyChanges(translation, transcription, text, example string, tagIds []string) {
	t.tagIDs = tagIds
	t.transcription = transcription
	t.text = text
	t.translation = translation
	t.example = example
	t.updatedAt = time.Now()
}

func (t *Translation) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            t.id,
		"authorID":      t.authorID,
		"createdAt":     t.createdAt,
		"updatedAt":     t.updatedAt,
		"translation":   t.translation,
		"transcription": t.transcription,
		"text":          t.text,
		"example":       t.example,
		"tagIDs":        t.tagIDs,
	}
}

func UnmarshalFromDB(
	id string,
	authorID string,
	createdAt time.Time,
	updatedAt time.Time,
	transcription string,
	translation string,
	text string,
	example string,
	tagIDs []string,
) *Translation {
	return &Translation{
		id:            id,
		authorID:      authorID,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		transcription: transcription,
		translation:   translation,
		text:          text,
		example:       example,
		tagIDs:        tagIDs,
	}
}
