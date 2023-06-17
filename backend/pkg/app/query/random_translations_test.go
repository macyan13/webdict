package query

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomTranslationsHandler_processLimit(t *testing.T) {
	type args struct {
		query RandomTranslations
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Limit is below min, default value is expecting",
			args{query: RandomTranslations{Limit: -4}},
			10,
		},
		{
			"Limit is below min, default value is expecting",
			args{query: RandomTranslations{Limit: 0}},
			10,
		},
		{
			"Limit is higher max, default value is expecting",
			args{query: RandomTranslations{Limit: 101}},
			10,
		},
		{
			"Limit is in acceptable range",
			args{query: RandomTranslations{Limit: 54}},
			54,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := RandomTranslationsHandler{}
			assert.Equalf(t, tt.want, h.processLimit(tt.args.query), "processLimit(%v)", tt.args.query)
		})
	}
}

func TestRandomTranslationsHandler_Handle(t *testing.T) {
	type fields struct {
		translationRepo TranslationViewRepository
	}
	type args struct {
		query RandomTranslations
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     RandomViews
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on getting random views from db",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetRandomViews", "authorID", "EN", []string{}, 10).Return(RandomViews{}, fmt.Errorf("error"))
				return fields{translationRepo: &repo}
			},
			args{RandomTranslations{AuthorID: "authorID", LangID: "EN", TagIds: []string{}, Limit: 10}},
			RandomViews{},
			assert.Error,
		},
		{
			"Positive case",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetRandomViews", "authorID", "EN", []string{}, 10).Return(
					RandomViews{
						Views: []TranslationView{{
							ID:            "testID",
							Source:        "TestText",
							Transcription: "testTranscription",
							Target:        `<a href=\"javascript:alert('XSS1')\" onmouseover=\"alert('XSS2')\"><br>TestMeaning</br><a>`,
							Example:       "testExample",
						}},
					}, nil)
				return fields{translationRepo: &repo}
			},
			args{RandomTranslations{AuthorID: "authorID", LangID: "EN", TagIds: []string{}, Limit: 10}},
			RandomViews{Views: []TranslationView{{
				ID:            "testID",
				Source:        "TestText",
				Transcription: "testTranscription",
				Target:        "<br>TestMeaning</br>",
				Example:       "testExample",
			}}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewRandomTranslationsHandler(f.translationRepo)
			got, err := h.Handle(tt.args.query)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.query)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v)", tt.args.query)
		})
	}
}
