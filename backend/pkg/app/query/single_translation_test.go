package query

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSingleTranslationHandler_Handle(t *testing.T) {
	type fields struct {
		translationRepo TranslationViewRepository
	}
	type args struct {
		cmd SingleTranslation
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     TranslationView
		wantErr  bool
	}{
		{
			"Error on query validation",
			func() fields {
				return fields{translationRepo: &MockTranslationViewRepository{}}
			},
			args{cmd: SingleTranslation{AuthorID: "testAuthor"}},
			TranslationView{},
			true,
		},
		{
			"Error on DB query",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetView", "trID", "testAuthor").Return(TranslationView{}, errors.New("testErr"))
				return fields{translationRepo: &repo}
			},
			args{cmd: SingleTranslation{ID: "trID", AuthorID: "testAuthor"}},
			TranslationView{},
			true,
		},
		{
			"Positive case",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetView", "trID", "testAuthor").Return(TranslationView{
					Target: "testTranslation",
				}, nil)
				return fields{translationRepo: &repo}
			},
			args{cmd: SingleTranslation{ID: "trID", AuthorID: "testAuthor"}},
			TranslationView{
				Target: "testTranslation",
			},
			false,
		},
		{
			"Check sanitization",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetView", "trID", "testAuthor").Return(TranslationView{
					Target: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">Test Target<a>`,
				}, nil)
				return fields{translationRepo: &repo}
			},
			args{cmd: SingleTranslation{ID: "trID", AuthorID: "testAuthor"}},
			TranslationView{
				Target: "Test Target",
			},
			false,
		},
	}

	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSingleTranslationHandler(tt.fieldsFn().translationRepo, v)
			got, err := h.Handle(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleTranslationValidation(t *testing.T) {
	type args struct {
		query SingleTranslation
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid case",
			args: args{
				query: SingleTranslation{
					ID:       "123",
					AuthorID: "456",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Missing ID",
			args: args{
				query: SingleTranslation{
					AuthorID: "456",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Missing AuthorID",
			args: args{
				query: SingleTranslation{
					ID: "123",
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
