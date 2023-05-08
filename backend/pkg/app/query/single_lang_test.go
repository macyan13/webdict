package query

import (
	"errors"
	"reflect"
	"testing"
)

func TestSingleLangHandler_Handle(t *testing.T) {
	type fields struct {
		langRepo LangViewRepository
	}
	type args struct {
		cmd SingleLang
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     LangView
		wantErr  bool
	}{
		{
			"Error on DB query",
			func() fields {
				repo := MockLangViewRepository{}
				repo.On("GetView", "langID", "testAuthor").Return(LangView{}, errors.New("testErr"))
				return fields{langRepo: &repo}
			},
			args{cmd: SingleLang{ID: "langID", AuthorID: "testAuthor"}},
			LangView{},
			true,
		},
		{
			"Positive case",
			func() fields {
				repo := MockLangViewRepository{}
				repo.On("GetView", "langID", "testAuthor").Return(LangView{
					ID:   "testID",
					Name: "en",
				}, nil)
				return fields{langRepo: &repo}
			},
			args{cmd: SingleLang{ID: "langID", AuthorID: "testAuthor"}},
			LangView{
				ID:   "testID",
				Name: "en",
			},
			false,
		},
		{
			"Case 3: check sanitization",
			func() fields {
				repo := MockLangViewRepository{}
				repo.On("GetView", "langID", "testAuthor").Return(LangView{
					ID:   "langID",
					Name: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">DE<a>`,
				}, nil)
				return fields{langRepo: &repo}
			},
			args{cmd: SingleLang{ID: "langID", AuthorID: "testAuthor"}},
			LangView{
				ID:   "langID",
				Name: "DE",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSingleLangHandler(tt.fieldsFn().langRepo)
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
