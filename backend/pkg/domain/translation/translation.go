package translation

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"time"
	"unicode/utf8"
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

func NewTranslation(text, transcription, translation, authorID, example string, tagIDs []string) (*Translation, error) {
	now := time.Now()
	tr := Translation{
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

	if err := tr.validate(); err != nil {
		return nil, err
	}

	return &tr, nil
}

func (t *Translation) ID() string {
	return t.id
}

func (t *Translation) AuthorID() string {
	return t.authorID
}

func (t *Translation) ApplyChanges(text, transcription, translation, example string, tagIds []string) error {
	updated := *t
	updated.applyChanges(text, transcription, translation, example, tagIds)

	if err := updated.validate(); err != nil {
		return err
	}

	t.applyChanges(text, transcription, translation, example, tagIds)
	return nil
}

func (t *Translation) applyChanges(text, transcription, translation, example string, tagIds []string) {
	t.tagIDs = tagIds
	t.transcription = transcription
	t.text = text
	t.translation = translation
	t.example = example
	t.updatedAt = time.Now()
}

func (t *Translation) validate() error {
	var result error
	if t.text == "" {
		result = multierror.Append(result, errors.New("text can not be empty"))
	}

	textCount := utf8.RuneCountInString(t.text)
	if textCount > 255 {
		result = multierror.Append(result, fmt.Errorf("text max size is 255 characters, %d passed (%s)", textCount, t.text))
	}

	transcriptionCount := utf8.RuneCountInString(t.transcription)
	if transcriptionCount > 255 {
		result = multierror.Append(result, fmt.Errorf("transcription max size is 255 characters, %d passed (%s)", transcriptionCount, t.transcription))
	}

	if t.translation == "" {
		result = multierror.Append(result, fmt.Errorf("translation can not be empty"))
	}

	translationCount := utf8.RuneCountInString(t.translation)
	if translationCount > 255 {
		result = multierror.Append(result, fmt.Errorf("translation max size is 255 characters, %d passed (%s)", translationCount, t.translation))
	}

	exampleCount := utf8.RuneCountInString(t.example)
	if utf8.RuneCountInString(t.example) > 255 {
		result = multierror.Append(result, fmt.Errorf("example max size is 255 characters, %d passed (%s)", exampleCount, t.example))
	}

	if t.authorID == "" {
		result = multierror.Append(result, fmt.Errorf("authorID can not be empty"))
	}

	tagsCount := len(t.tagIDs)
	if tagsCount > 5 {
		result = multierror.Append(result, fmt.Errorf("tag max amount is 5, %d passed", tagsCount))
	}

	return result
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
