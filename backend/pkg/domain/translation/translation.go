package translation

import (
	"github.com/google/uuid"
	"time"
)

// todo: clean up unused getter after read DB repository implementation
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

func (t *Translation) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Translation) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Translation) Transcription() string {
	return t.transcription
}

func (t *Translation) Translation() string {
	return t.translation
}

func (t *Translation) Text() string {
	return t.text
}

func (t *Translation) Example() string {
	return t.example
}

func (t *Translation) TagIds() []string {
	return t.tagIds
}

func (t *Translation) ApplyChanges(translation, transcription, text, example string, tagIds []string) {
	t.tagIds = tagIds
	t.transcription = transcription
	t.text = text
	t.translation = translation
	t.example = example
	t.updatedAt = time.Now()
}

func (t *Translation) IsAuthor(authorId string) bool {
	return t.authorId == authorId
}
