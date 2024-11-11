package query

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLastTranslationsHandler_Handle(t *testing.T) {
	type fields struct {
		translationRepo TranslationViewRepository
	}
	type args struct {
		query SearchTranslations
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     LastTranslationViews
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error validation",
			func() fields {
				return fields{translationRepo: &MockTranslationViewRepository{}}
			},
			args{SearchTranslations{AuthorID: "authorID", PageSize: 10, Page: 1}},
			LastTranslationViews{},
			assert.Error,
		},
		{
			"Error on getting last views from repository",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 1, []string{}).Return(LastTranslationViews{}, fmt.Errorf("error"))
				return fields{translationRepo: &repo}
			},
			args{SearchTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 10, Page: 1, TagIds: []string{}}},
			LastTranslationViews{},
			assert.Error,
		},
		{
			"Search by tags",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 1, []string{"tag1", "tag2"}).Return(
					LastTranslationViews{
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
			args{SearchTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 10, Page: 1, TagIds: []string{"tag1", "tag2"}}},
			LastTranslationViews{Views: []TranslationView{{
				ID:            "testID",
				Source:        "TestText",
				Transcription: "testTranscription",
				Target:        "<br>TestMeaning</br>",
				Example:       "testExample",
			}}},
			assert.NoError,
		},
		{
			"Search by tags",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 1, []string{"tag1", "tag2"}).Return(
					LastTranslationViews{
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
			args{SearchTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 10, Page: 1, TagIds: []string{"tag1", "tag2"}}},
			LastTranslationViews{Views: []TranslationView{{
				ID:            "testID",
				Source:        "TestText",
				Transcription: "testTranscription",
				Target:        "<br>TestMeaning</br>",
				Example:       "testExample",
			}}},
			assert.NoError,
		},
		{
			"Search by target part",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViewsByTargetPart", "authorID", "EN", "targetPart", 10, 1).Return(
					LastTranslationViews{
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
			args{SearchTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 10, Page: 1, TargetPart: "targetPart"}},
			LastTranslationViews{Views: []TranslationView{{
				ID:            "testID",
				Source:        "TestText",
				Transcription: "testTranscription",
				Target:        "<br>TestMeaning</br>",
				Example:       "testExample",
			}}},
			assert.NoError,
		},
		{
			"Search by source part",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetLastViewsBySourcePart", "authorID", "EN", "sourcePart", 10, 1).Return(
					LastTranslationViews{
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
			args{SearchTranslations{AuthorID: "authorID", LangID: "EN", PageSize: 10, Page: 1, SourcePart: "sourcePart"}},
			LastTranslationViews{Views: []TranslationView{{
				ID:            "testID",
				Source:        "TestText",
				Transcription: "testTranscription",
				Target:        "<br>TestMeaning</br>",
				Example:       "testExample",
			}}},
			assert.NoError,
		},
	}

	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewSearchTranslationsHandler(f.translationRepo, v)
			got, err := h.Handle(tt.args.query)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.query)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v)", tt.args.query)
		})
	}
}

func TestSearchTranslationsValidation(t *testing.T) {
	type args struct {
		query SearchTranslations
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Error on missing AuthorID",
			args{SearchTranslations{LangID: "EN", PageSize: 10, Page: 1}},
			assert.Error,
		},
		{
			"Error on missing LangID",
			args{SearchTranslations{AuthorID: "123", PageSize: 10, Page: 1}},
			assert.Error,
		},
		{
			"Error on missing PageSize",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", Page: 1}},
			assert.Error,
		},
		{
			"Error on missing Page",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 10}},
			assert.Error,
		},
		{
			"Error on conflicting SourcePart and TargetPart",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 10, Page: 1, SourcePart: "source", TargetPart: "target"}},
			assert.Error,
		},
		{
			"Error on conflicting TagIDs and SourcePart",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 10, Page: 1, TagIds: []string{"tag1", "tag2"}, SourcePart: "source"}},
			assert.Error,
		},
		{
			"Error on invalid PageSize (less than 1)",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 0, Page: 1}},
			assert.Error,
		},
		{
			"Error on invalid PageSize (greater than 200)",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 201, Page: 1}},
			assert.Error,
		},
		{
			"Error on invalid Page (less than 1)",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 10, Page: 0}},
			assert.Error,
		},
		{
			"Valid input",
			args{SearchTranslations{AuthorID: "123", LangID: "EN", PageSize: 10, Page: 1}},
			assert.NoError,
		},
	}

	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.args.query)
			tt.wantErr(t, err, fmt.Sprintf("v.Struct(%v)", tt.args.query))
		})
	}
}
