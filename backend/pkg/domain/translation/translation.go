package translation

import (
	"github.com/google/uuid"
	"time"
)

type Lang string

const EN Lang = "en"

type Translation struct {
	id            string
	text          string
	transcription string
	translation   string
	authorID      string
	example       string
	tagIDs        []string
	createdAt     time.Time
	updatedAt     time.Time
	lang          Lang
}

func NewTranslation(text, transcription, translation, authorID, example string, tagIDs []string) *Translation {
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
		lang:          EN,
	}
}

func (t *Translation) ID() string {
	return t.id
}

func (t *Translation) AuthorID() string {
	return t.authorID
}

func (t *Translation) ApplyChanges(text, transcription, translation, example string, tagIds []string) {
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
		"text":          t.text,
		"transcription": t.transcription,
		"translation":   t.translation,
		"authorID":      t.authorID,
		"example":       t.example,
		"tagIDs":        t.tagIDs,
		"createdAt":     t.createdAt,
		"updatedAt":     t.updatedAt,
		"lang":          t.lang,
	}
}

func UnmarshalFromDB(
	id string,
	text string,
	transcription string,
	translation string,
	authorID string,
	example string,
	tagIDs []string,
	createdAt time.Time,
	updatedAt time.Time,
	lang Lang,
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
		lang:          lang,
	}
}
