package translation

import (
	"testing"
	"time"
)

import "github.com/stretchr/testify/assert"

func TestNewTranslation(t *testing.T) {
	translation := NewTranslation("", "", "", "", "", []string{})
	assert.Equal(t, translation.updatedAt, translation.createdAt, "NewTranslation - createdAt and updatedAt are the same")
}

func TestTranslation_ApplyChanges(t *testing.T) {
	tr := NewTranslation("", "", "", "", "", []string{})
	translation := "test"
	transcription := "[test]"
	text := "text"
	example := "exampleTest"
	tg := "testTag"
	updatedAt := tr.updatedAt

	time.Sleep(time.Second)
	tr.ApplyChanges(text, transcription, translation, example, []string{tg})

	assert.Equal(t, tr.translation, translation)
	assert.Equal(t, tr.text, text)
	assert.Equal(t, tr.example, example)
	assert.Equal(t, tr.translation, translation)
	assert.Greaterf(t, tr.updatedAt, updatedAt, "Tag.ApplyChanges - updatedAt should be greater createdAt")
	assert.Equal(t, tr.tagIDs[0], tg)
}

func TestUnmarshalFromDB(t *testing.T) {
	translation := Translation{
		id:            "testId",
		authorID:      "testAuthor",
		createdAt:     time.Now().Add(5 * time.Second),
		updatedAt:     time.Now().Add(10 * time.Second),
		transcription: "testTranscription",
		translation:   "testTranslation",
		text:          "testText",
		example:       "testExample",
		tagIDs:        []string{"tag1", "tag2"},
		lang:          EN,
	}

	assert.Equal(t, &translation, UnmarshalFromDB(
		translation.id,
		translation.text,
		translation.transcription,
		translation.translation,
		translation.authorID,
		translation.example,
		translation.tagIDs,
		translation.createdAt,
		translation.updatedAt,
		EN,
	))
}
