package query

import (
	"errors"
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
			"Case 1: error on DB query",
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
			"Case 2: positive case",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetView", "trID", "testAuthor").Return(TranslationView{
					Translation: "testTranslation",
				}, nil)
				return fields{translationRepo: &repo}
			},
			args{cmd: SingleTranslation{ID: "trID", AuthorID: "testAuthor"}},
			TranslationView{
				Translation: "testTranslation",
			},
			false,
		},
		{
			"Case 3: check sanitization",
			func() fields {
				repo := MockTranslationViewRepository{}
				repo.On("GetView", "trID", "testAuthor").Return(TranslationView{
					Translation: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">Test Translation<a>`,
				}, nil)
				return fields{translationRepo: &repo}
			},
			args{cmd: SingleTranslation{ID: "trID", AuthorID: "testAuthor"}},
			TranslationView{
				Translation: "Test Translation",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSingleTranslationHandler(tt.fieldsFn().translationRepo)
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
