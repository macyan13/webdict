package translation

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
	"unicode/utf8"
)

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
	langID        string
}

func NewTranslation(source, transcription, target, authorID, example string, tagIDs []string, langID string) (*Translation, error) {
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
		langID:        langID,
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

func (t *Translation) LangID() string {
	return t.langID
}

func (t *Translation) ApplyChanges(source, transcription, target, example string, tagIDs []string, langID string) error {
	updated := *t
	updated.applyChanges(source, transcription, target, example, tagIDs, langID)

	if err := updated.validate(); err != nil {
		return err
	}

	t.applyChanges(source, transcription, target, example, tagIDs, langID)
	return nil
}

func (t *Translation) applyChanges(source, transcription, target, example string, tagIDs []string, langID string) {
	t.tagIDs = tagIDs
	t.transcription = transcription
	t.source = source
	t.target = target
	t.example = example
	t.updatedAt = time.Now()
	t.langID = langID
}

func (t *Translation) validate() error {
	var err error
	if t.source == "" {
		err = errors.Join(errors.New("source can not be empty"), err)
	}

	textCount := utf8.RuneCountInString(t.source)
	if textCount > 255 {
		err = errors.Join(fmt.Errorf("source max size is 255 characters, %d passed (%s)", textCount, t.source), err)
	}

	transcriptionCount := utf8.RuneCountInString(t.transcription)
	if transcriptionCount > 255 {
		err = errors.Join(fmt.Errorf("transcription max size is 255 characters, %d passed (%s)", transcriptionCount, t.transcription), err)
	}

	if t.target == "" {
		err = errors.Join(fmt.Errorf("target can not be empty"), err)
	}

	translationCount := utf8.RuneCountInString(t.target)
	if translationCount > 255 {
		err = errors.Join(fmt.Errorf("target max size is 255 characters, %d passed (%s)", translationCount, t.target), err)
	}

	exampleCount := utf8.RuneCountInString(t.example)
	if utf8.RuneCountInString(t.example) > 255 {
		err = errors.Join(fmt.Errorf("example max size is 255 characters, %d passed (%s)", exampleCount, t.example), err)
	}

	if t.authorID == "" {
		err = errors.Join(fmt.Errorf("authorID can not be empty"), err)
	}

	tagsCount := len(t.tagIDs)
	if tagsCount > 5 {
		err = errors.Join(fmt.Errorf("tag max amount is 5, %d passed", tagsCount), err)
	}

	if t.langID == "" {
		err = errors.Join(fmt.Errorf("langID can not be empty"), err)
	}

	return err
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
		"langID":        t.langID,
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
	langID string,
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
		langID:        langID,
	}
}
