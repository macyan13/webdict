package query

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			"Error on query validation",
			func() fields {
				return fields{translationRepo: &MockTranslationViewRepository{}}
			},
			args{RandomTranslations{AuthorID: "authorID", LangID: "EN", TagIds: []string{}, Limit: 0}},
			RandomViews{},
			assert.Error,
		},
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
	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewRandomTranslationsHandler(f.translationRepo, v)
			got, err := h.Handle(tt.args.query)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.query)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v)", tt.args.query)
		})
	}
}

func TestRandomTranslationsValidation(t *testing.T) {
	type args struct {
		query RandomTranslations
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid case",
			args: args{
				query: RandomTranslations{
					AuthorID: "123",
					LangID:   "en",
					TagIds:   []string{"tag1", "tag2"},
					Limit:    100,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Missing AuthorID",
			args: args{
				query: RandomTranslations{
					LangID: "en",
					TagIds: []string{"tag1", "tag2"},
					Limit:  100,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Missing LangID",
			args: args{
				query: RandomTranslations{
					AuthorID: "123",
					TagIds:   []string{"tag1", "tag2"},
					Limit:    100,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid Limit (less than 1)",
			args: args{
				query: RandomTranslations{
					AuthorID: "123",
					LangID:   "en",
					TagIds:   []string{"tag1", "tag2"},
					Limit:    0,
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Invalid Limit (greater than 200)",
			args: args{
				query: RandomTranslations{
					AuthorID: "123",
					LangID:   "en",
					TagIds:   []string{"tag1", "tag2"},
					Limit:    201,
				},
			},
			wantErr: assert.Error,
		},
	}

	v := validator.New()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.args.query)
			test.wantErr(t, err)
		})
	}
}
