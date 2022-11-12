package translation

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"testing"
	"time"
)

import "github.com/stretchr/testify/assert"

func TestNewTranslation(t *testing.T) {
	translation := NewTranslation(Data{})
	assert.Equal(t, translation.updatedAt, translation.createdAt)
}

func TestTranslation_ApplyChanges(t *testing.T) {
	tr := NewTranslation(Data{})
	updatedAt := tr.updatedAt

	translation := "test"
	transcription := "[test]"
	text := "text"
	example := "exampleTest"
	tg := "testTag"

	data := Data{
		Request: Request{
			Transcription: transcription,
			Translation:   translation,
			Text:          text,
			Example:       example,
		},
		Tags: []*tag.Tag{{
			tag: tg,
		}},
	}

	time.Sleep(time.Second)

	tr.applyChanges(data)

	assert := assert.New(t)
	assert.Equal(tr.translation, translation)
	assert.Equal(tr.text, text)
	assert.Equal(tr.example, example)
	assert.Equal(tr.translation, translation)
	assert.Greaterf(tr.updatedAt, updatedAt, "error message %s", "formatted")
	assert.Equal(tr.tagIds[0].Tag, tg)
}
