package query

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
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
			"Error on query validation",
			func() fields {
				return fields{langRepo: &MockLangViewRepository{}}
			},
			args{cmd: SingleLang{ID: "", AuthorID: "testAuthor"}},
			LangView{},
			true,
		},
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

	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSingleLangHandler(tt.fieldsFn().langRepo, v)
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

func TestSingleLangValidation(t *testing.T) {
	type args struct {
		query SingleLang
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid case",
			args: args{
				query: SingleLang{
					ID:       "123",
					AuthorID: "456",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Missing ID",
			args: args{
				query: SingleLang{
					AuthorID: "456",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Missing AuthorID",
			args: args{
				query: SingleLang{
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
