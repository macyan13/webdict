package translation

import (
	"strings"
	"testing"
	"time"
)

import "github.com/stretchr/testify/assert"

func TestNewTranslation_PositiveCase(t *testing.T) {
	translation, err := NewTranslation("new", "", "new", "new", "", []string{}, "EN")
	assert.Nil(t, err)
	assert.Equal(t, translation.updatedAt, translation.createdAt, "NewTranslation - createdAt and updatedAt are the same")
}

func TestNewTranslation_ValidationError(t *testing.T) {
	tr, err := NewTranslation("", "", "new", "new", "", []string{}, "EN")
	assert.Nil(t, tr)
	assert.True(t, strings.Contains(err.Error(), "source can not be empty"))
}

func TestTranslation_ApplyChanges_PositiveCase(t *testing.T) {
	tr, err := NewTranslation("new", "new", "new", "new", "new", []string{}, "EN")
	assert.Nil(t, err)
	translation := "test"
	transcription := "[test]"
	text := "source"
	example := "exampleTest"
	tg := "testTag"
	updatedAt := tr.updatedAt
	langID := "lang1"

	time.Sleep(time.Second)
	err = tr.ApplyChanges(text, transcription, translation, example, []string{tg}, langID)
	assert.Nil(t, err)

	assert.Equal(t, tr.target, translation)
	assert.Equal(t, tr.source, text)
	assert.Equal(t, tr.example, example)
	assert.Equal(t, tr.target, translation)
	assert.Greaterf(t, tr.updatedAt, updatedAt, "Name.ApplyChanges - updatedAt should be greater createdAt")
	assert.Equal(t, tr.tagIDs[0], tg)
	assert.Equal(t, langID, tr.langID)
}

func TestTranslation_ApplyChanges_ValidationError(t *testing.T) {
	tr, err := NewTranslation("new", "new", "new", "new", "new", []string{}, "EN")
	assert.Nil(t, err)

	err = tr.ApplyChanges("", "", "test", "", []string{}, "DE")
	assert.True(t, strings.Contains(err.Error(), "source can not be empty"))
	assert.Equal(t, "new", tr.target)
}

func TestUnmarshalFromDB(t *testing.T) {
	translation := Translation{
		id:            "testId",
		authorID:      "testAuthor",
		createdAt:     time.Now().Add(5 * time.Second),
		updatedAt:     time.Now().Add(10 * time.Second),
		transcription: "testTranscription",
		target:        "testTranslation",
		source:        "testText",
		example:       "testExample",
		tagIDs:        []string{"tag1", "tag2"},
		langID:        "EN",
	}

	assert.Equal(t, &translation, UnmarshalFromDB(
		translation.id,
		translation.source,
		translation.transcription,
		translation.target,
		translation.authorID,
		translation.example,
		translation.tagIDs,
		translation.createdAt,
		translation.updatedAt,
		"EN",
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
		langID        string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Source is empty",
			fields{
				translation: "test",
				authorID:    "test",
				langID:      "lang1",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "source can not be empty"), i)
				return true
			},
		},
		{
			"Source is too long",
			fields{
				text:        string(make([]rune, 256)),
				translation: "test",
				authorID:    "test",
				langID:      "lang1",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "source max size is 255 characters, 256 passed ("), i)
				return true
			},
		},
		{
			"Target is too long",
			fields{
				text:          "test",
				transcription: string(make([]rune, 256)),
				translation:   "test",
				authorID:      "test",
				langID:        "EN",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "transcription max size is 255 characters, 256 passed"), i)
				return true
			},
		},
		{
			"Target is empty",
			fields{
				text:     "test",
				authorID: "test",
				langID:   "lang1",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "target can not be empty"), i)
				return true
			},
		},
		{
			"Target is too long",
			fields{
				text:        "test",
				translation: string(make([]rune, 256)),
				authorID:    "test",
				langID:      "lang1",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "target max size is 255 characters, 256 passed"), i)
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
				langID:      "lang1",
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
				langID:      "lang1",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "authorID can not be empty"), i)
				return true
			},
		},
		{
			"LangID is missed",
			fields{
				text:        "test",
				translation: "test",
				authorID:    "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "langID can not be empty"), i)
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
				langID:      "lang1",
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
				langID:   "lang1",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "source can not be empty"), i)
				assert.True(t, strings.Contains(err.Error(), "target can not be empty"), i)
				return true
			},
		},
		{
			"Positive case",
			fields{
				text:        "test",
				translation: "test",
				authorID:    "test",
				langID:      "lang1",
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
				source:        tt.fields.text,
				transcription: tt.fields.transcription,
				target:        tt.fields.translation,
				authorID:      tt.fields.authorID,
				example:       tt.fields.example,
				tagIDs:        tt.fields.tagIDs,
				langID:        tt.fields.langID,
			}
			tt.wantErr(t, tr.validate(), "validate()")
		})
	}
}
