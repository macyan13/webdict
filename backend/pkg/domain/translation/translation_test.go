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
	tr.ApplyChanges(translation, transcription, text, example, []string{tg})

	assert := assert.New(t)
	assert.Equal(tr.translation, translation)
	assert.Equal(tr.text, text)
	assert.Equal(tr.example, example)
	assert.Equal(tr.translation, translation)
	assert.Greaterf(tr.updatedAt, updatedAt, "Tag.ApplyChanges - updatedAt should be greater createdAt")
	assert.Equal(tr.tagIds[0], tg)
}
