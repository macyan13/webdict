package translation

import (
	"strings"
	"testing"
	"time"
)

import "github.com/stretchr/testify/assert"

func TestNewTranslation_PositiveCase(t *testing.T) {
	translation, err := NewTranslation("new", "", "new", "new", "", []string{})
	assert.Nil(t, err)
	assert.Equal(t, translation.updatedAt, translation.createdAt, "NewTranslation - createdAt and updatedAt are the same")
}

func TestNewTranslation_ValidationError(t *testing.T) {
	tr, err := NewTranslation("", "", "new", "new", "", []string{})
	assert.Nil(t, tr)
	assert.True(t, strings.Contains(err.Error(), "text can not be empty"))
}

func TestTranslation_ApplyChanges_PositiveCase(t *testing.T) {
	tr, err := NewTranslation("new", "new", "new", "new", "new", []string{})
	assert.Nil(t, err)
	translation := "test"
	transcription := "[test]"
	text := "text"
	example := "exampleTest"
	tg := "testTag"
	updatedAt := tr.updatedAt

	time.Sleep(time.Second)
	err = tr.ApplyChanges(text, transcription, translation, example, []string{tg})
	assert.Nil(t, err)

	assert.Equal(t, tr.meaning, translation)
	assert.Equal(t, tr.text, text)
	assert.Equal(t, tr.example, example)
	assert.Equal(t, tr.meaning, translation)
	assert.Greaterf(t, tr.updatedAt, updatedAt, "Tag.ApplyChanges - updatedAt should be greater createdAt")
	assert.Equal(t, tr.tagIDs[0], tg)
}

func TestTranslation_ApplyChanges_ValidationError(t *testing.T) {
	tr, err := NewTranslation("new", "new", "new", "new", "new", []string{})
	assert.Nil(t, err)

	err = tr.ApplyChanges("", "", "test", "", []string{})
	assert.True(t, strings.Contains(err.Error(), "text can not be empty"))
	assert.Equal(t, "new", tr.meaning)
}

func TestUnmarshalFromDB(t *testing.T) {
	translation := Translation{
		id:            "testId",
		authorID:      "testAuthor",
		createdAt:     time.Now().Add(5 * time.Second),
		updatedAt:     time.Now().Add(10 * time.Second),
		transcription: "testTranscription",
		meaning:       "testTranslation",
		text:          "testText",
		example:       "testExample",
		tagIDs:        []string{"tag1", "tag2"},
		lang:          EN,
	}

	assert.Equal(t, &translation, UnmarshalFromDB(
		translation.id,
		translation.text,
		translation.transcription,
		translation.meaning,
		translation.authorID,
		translation.example,
		translation.tagIDs,
		translation.createdAt,
		translation.updatedAt,
		EN,
	))
}

func TestTranslation_validate(t *testing.T) {
	type fields struct {
		text          string
		transcription string
		translation   string
		authorID      string
		example       string
		tagIDs        []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Text is empty",
			fields{
				translation: "test",
				authorID:    "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "text can not be empty"), i)
				return true
			},
		},
		{
			"Text is too long",
			fields{
				text:        string(make([]rune, 256)),
				translation: "test",
				authorID:    "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "text max size is 255 characters, 256 passed ("), i)
				return true
			},
		},
		{
			"Meaning is too long",
			fields{
				text:          "test",
				transcription: string(make([]rune, 256)),
				translation:   "test",
				authorID:      "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "transcription max size is 255 characters, 256 passed"), i)
				return true
			},
		},
		{
			"Meaning is empty",
			fields{
				text:     "test",
				authorID: "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "meaning can not be empty"), i)
				return true
			},
		},
		{
			"Meaning is too long",
			fields{
				text:        "test",
				translation: string(make([]rune, 256)),
				authorID:    "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "meaning max size is 255 characters, 256 passed"), i)
				return true
			},
		},
		{
			"Example is too long",
			fields{
				text:        "test",
				translation: "test",
				example:     string(make([]rune, 256)),
				authorID:    "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "example max size is 255 characters, 256 passed"), i)
				return true
			},
		},
		{
			"authorID is too long",
			fields{
				text:        "test",
				translation: "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "authorID can not be empty"), i)
				return true
			},
		},
		{
			"Too many tags",
			fields{
				text:        "test",
				translation: "test",
				authorID:    "test",
				tagIDs:      []string{"test", "test", "test", "test", "test", "test"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "tag max amount is 5, 6 passed"), i)
				return true
			},
		},
		{
			"Multiple errors",
			fields{
				authorID: "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "text can not be empty"), i)
				assert.True(t, strings.Contains(err.Error(), "meaning can not be empty"), i)
				return true
			},
		},
		{
			"Positive case",
			fields{
				text:        "test",
				translation: "test",
				authorID:    "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Translation{
				text:          tt.fields.text,
				transcription: tt.fields.transcription,
				meaning:       tt.fields.translation,
				authorID:      tt.fields.authorID,
				example:       tt.fields.example,
				tagIDs:        tt.fields.tagIDs,
			}
			tt.wantErr(t, tr.validate(), "validate()")
		})
	}
}
