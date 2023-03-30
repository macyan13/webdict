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
	source        string
	transcription string
	target        string
	authorID      string
	example       string
	tagIDs        []string
	createdAt     time.Time
	updatedAt     time.Time
	lang          Lang
}

func NewTranslation(source, transcription, target, authorID, example string, tagIDs []string) (*Translation, error) {
	now := time.Now()
	tr := Translation{
		id:            uuid.New().String(),
		authorID:      authorID,
		createdAt:     now,
		updatedAt:     now,
		target:        target,
		transcription: transcription,
		source:        source,
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

func (t *Translation) ApplyChanges(source, transcription, target, example string, tagIds []string) error {
	updated := *t
	updated.applyChanges(source, transcription, target, example, tagIds)

	if err := updated.validate(); err != nil {
		return err
	}

	t.applyChanges(source, transcription, target, example, tagIds)
	return nil
}

func (t *Translation) applyChanges(source, transcription, target, example string, tagIds []string) {
	t.tagIDs = tagIds
	t.transcription = transcription
	t.source = source
	t.target = target
	t.example = example
	t.updatedAt = time.Now()
}

func (t *Translation) validate() error {
	var result error
	if t.source == "" {
		result = multierror.Append(result, errors.New("source can not be empty"))
	}

	textCount := utf8.RuneCountInString(t.source)
	if textCount > 255 {
		result = multierror.Append(result, fmt.Errorf("source max size is 255 characters, %d passed (%s)", textCount, t.source))
	}

	transcriptionCount := utf8.RuneCountInString(t.transcription)
	if transcriptionCount > 255 {
		result = multierror.Append(result, fmt.Errorf("transcription max size is 255 characters, %d passed (%s)", transcriptionCount, t.transcription))
	}

	if t.target == "" {
		result = multierror.Append(result, fmt.Errorf("target can not be empty"))
	}

	translationCount := utf8.RuneCountInString(t.target)
	if translationCount > 255 {
		result = multierror.Append(result, fmt.Errorf("target max size is 255 characters, %d passed (%s)", translationCount, t.target))
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
		"source":        t.source,
		"transcription": t.transcription,
		"target":        t.target,
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
	source string,
	transcription string,
	target string,
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
		target:        target,
		source:        source,
		example:       example,
		tagIDs:        tagIDs,
		lang:          lang,
	}
}
