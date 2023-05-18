package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTagView_sanitize(t *testing.T) {
	tests := []struct {
		name         string
		rawTag       string
		sanitizedTag string
	}{
		{
			"Case 1: tag without malicious content",
			"cleanTag",
			"cleanTag",
		},
		{
			"Case 1: tag with malicious content",
			`<a onblur="alert(secret)" href="http://www.test.com">Test</a>`,
			`Test`,
		},
	}
	id := "<must_not_change>"
	sanitizer := newStrictSanitizer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &TagView{
				ID:  id,
				Tag: tt.rawTag,
			}
			v.sanitize(sanitizer)
			assert.Equal(t, tt.sanitizedTag, v.Tag)
			assert.Equal(t, id, v.ID)
		})
	}
}

func TestTranslationView_sanitize(t *testing.T) {
	type fields struct {
		Text          string
		Transcription string
		Translation   string
		Example       string
		Tag           string
		LangView      LangView
	}
	tests := []struct {
		name            string
		rawFields       fields
		sanitizedFields fields
	}{
		{
			"Case 1: translation without malicious content",
			fields{
				Text:          "testText",
				Transcription: "[testTranscription]",
				Translation:   "testTranslation",
				Example:       "testExample",
				Tag:           "testTag",
				LangView: LangView{
					ID:   "test",
					Name: "English",
				},
			},
			fields{
				Text:          "testText",
				Transcription: "[testTranscription]",
				Translation:   "testTranslation",
				Example:       "testExample",
				Tag:           "testTag",
				LangView: LangView{
					ID:   "test",
					Name: "English",
				},
			},
		},
		{
			"Case 2: translation with malicious content",
			fields{
				Text:          `<a onblur="alert(secret)" href="http://www.test.com">Source</a>`,
				Transcription: `<a onblur="alert(secret)" href="http://www.test.com">[Target]</a>`,
				Translation:   `<a onblur="alert(secret)" href="http://www.test.com">Target</a>`,
				Example:       `<a onblur="alert(secret)" href="http://www.test.com">Example</a>`,
				Tag:           `<a onblur="alert(secret)" href="http://www.test.com">Tag</a>`,
				LangView: LangView{
					ID:   "langID",
					Name: `<a onblur="alert(secret)" href="http://www.test.com">English</a>`,
				},
			},
			fields{
				Text:          `&lt;a href=&#34;http://www.test.com&#34; rel=&#34;nofollow&#34;&gt;Source&lt;/a&gt;`,
				Transcription: `&lt;a href=&#34;http://www.test.com&#34; rel=&#34;nofollow&#34;&gt;[Target]&lt;/a&gt;`,
				Translation:   `<a href="http://www.test.com" rel="nofollow">Target</a>`,
				Example:       `<a href="http://www.test.com" rel="nofollow">Example</a>`,
				Tag:           `Tag`,
				LangView: LangView{
					ID:   "langID",
					Name: "English",
				},
			},
		},
	}
	id := "<must_not_change>"
	strictSntz := newStrictSanitizer()
	reachSntz := newRichTextSanitizer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &TranslationView{
				ID:            id,
				Source:        tt.rawFields.Text,
				Transcription: tt.rawFields.Transcription,
				Target:        tt.rawFields.Translation,
				Example:       tt.rawFields.Example,
				Tags:          []TagView{{Tag: tt.rawFields.Tag}},
				Lang:          tt.rawFields.LangView,
			}
			v.sanitize(strictSntz, reachSntz)
			assert.Equal(t, id, v.ID)
			assert.Equal(t, tt.sanitizedFields.Text, v.Source)
			assert.Equal(t, tt.sanitizedFields.Transcription, v.Transcription)
			assert.Equal(t, tt.sanitizedFields.Translation, v.Target)
			assert.Equal(t, tt.sanitizedFields.Example, v.Example)
			assert.Equal(t, tt.sanitizedFields.Tag, v.Tags[0].Tag)
			assert.Equal(t, tt.sanitizedFields.LangView, v.Lang)
		})
	}
}
