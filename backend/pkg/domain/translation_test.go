package domain

import (
	"testing"
	"time"
)

import "github.com/stretchr/testify/assert"

func TestNewTranslation(t *testing.T) {
	translation := NewTranslation(TranslationRequest{})
	assert.Equal(t, translation.UpdatedAt, translation.CreatedAt)
}

func TestTranslation_ApplyChanges(t *testing.T) {
	entity := NewTranslation(TranslationRequest{})
	updatedAt := entity.UpdatedAt

	translation := "test"
	transcription := "[test]"
	text := "text"
	example := "exampleTest"
	request := TranslationRequest{
		Transcription: transcription,
		Translation:   translation,
		Text:          text,
		Example:       example,
	}

	time.Sleep(time.Second)

	entity.ApplyChanges(request)

	assert := assert.New(t)
	assert.Equal(entity.Translation, translation)
	assert.Equal(entity.Text, text)
	assert.Equal(entity.Example, example)
	assert.Equal(entity.Translation, translation)
	assert.Greaterf(entity.UpdatedAt, updatedAt, "error message %s", "formatted")
}
