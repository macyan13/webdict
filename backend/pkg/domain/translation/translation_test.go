package translation

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"testing"
	"time"
)

import "github.com/stretchr/testify/assert"

func TestNewTranslation(t *testing.T) {
	translation := newTranslation(data{})
	assert.Equal(t, translation.UpdatedAt, translation.CreatedAt)
}

func TestTranslation_ApplyChanges(t *testing.T) {
	tr := newTranslation(data{})
	updatedAt := tr.UpdatedAt

	translation := "test"
	transcription := "[test]"
	text := "text"
	example := "exampleTest"
	tg := "testTag"

	data := data{
		Request: Request{
			Transcription: transcription,
			Translation:   translation,
			Text:          text,
			Example:       example,
		},
		Tags: []*tag.Tag{{
			Tag: tg,
		}},
	}

	time.Sleep(time.Second)

	tr.applyChanges(data)

	assert := assert.New(t)
	assert.Equal(tr.Translation, translation)
	assert.Equal(tr.Text, text)
	assert.Equal(tr.Example, example)
	assert.Equal(tr.Translation, translation)
	assert.Greaterf(tr.UpdatedAt, updatedAt, "error message %s", "formatted")
	assert.Equal(tr.Tags[0].Tag, tg)
}
