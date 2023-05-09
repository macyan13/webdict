package query

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLastTranslationsHandler_processParameters(t *testing.T) {
	type args struct {
		query LastTranslations
	}
	tests := []struct {
		name         string
		args         args
		wantPageSize int
		wantPage     int
	}{
		{
			"Page size is less then 1",
			args{LastTranslations{PageSize: -1, Page: 5}},
			10,
			5,
		},
		{
			"Page size is greater than 100",
			args{LastTranslations{PageSize: 101, Page: 5}},
			10,
			5,
		},
		{"Page is less then 1",
			args{LastTranslations{PageSize: 10, Page: -1}},
			10,
			1,
		},
		{
			"Positive case",
			args{LastTranslations{PageSize: 10, Page: 5}},
			10,
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := LastTranslationsHandler{}
			limit, page := h.processParameters(tt.args.query)
			assert.Equalf(t, tt.wantPageSize, limit, "processParameters(%v)", tt.args.query)
			assert.Equalf(t, tt.wantPage, page, "processParameters(%v)", tt.args.query)
		})
	}
}

func TestLastTranslationsHandler_Handle(t *testing.T) {
	type fields struct {
		translationRepo TranslationViewRepository
	}
	type args struct {
		query LastTranslations
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     LastViews
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on getting last views from repository",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViews", "authorID", "EN", 10, 1, []string{}).Return(LastViews{}, fmt.Errorf("error"))
				return fields{translationRepo: &repo}
			},
			args{LastTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 10, Page: 1, TagIds: []string{}}},
			LastViews{},
			assert.Error,
		},
		{
			"Positive case",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViews", "authorID", "EN", 10, 1, []string{}).Return(
					LastViews{
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
			args{LastTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 105, Page: 1, TagIds: []string{}}},
			LastViews{Views: []TranslationView{{
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
			h := NewLastTranslationsHandler(f.translationRepo)
			got, err := h.Handle(tt.args.query)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.query)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v)", tt.args.query)
		})
	}
}
