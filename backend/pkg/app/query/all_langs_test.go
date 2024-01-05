package query

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
)

func TestAllLangsHandler_Handle(t *testing.T) {
	type fields struct {
		langRepo LangViewRepository
	}
	type args struct {
		cmd AllLangs
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     []LangView
		wantErr  bool
	}{
		{
			"Error validation",
			func() fields {
				return fields{langRepo: &MockLangViewRepository{}}
			},
			args{cmd: AllLangs{AuthorID: ""}},
			nil,
			true,
		},
		{
			"Error on DB query",
			func() fields {
				repo := MockLangViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return(nil, errors.New("testErr"))
				return fields{langRepo: &repo}
			},
			args{cmd: AllLangs{AuthorID: "testAuthor"}},
			nil,
			true,
		},
		{
			"Positive case",
			func() fields {
				repo := MockLangViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return([]LangView{{
					ID:   "testId",
					Name: "en",
				}}, nil)
				return fields{langRepo: &repo}
			},
			args{cmd: AllLangs{AuthorID: "testAuthor"}},
			[]LangView{{
				ID:   "testId",
				Name: "en",
			}},
			false,
		},
		{
			"Check sanitization",
			func() fields {
				repo := MockLangViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return([]LangView{{
					ID:   "testId",
					Name: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')"><br>EN</br><a>`,
				}}, nil)
				return fields{langRepo: &repo}
			},
			args{cmd: AllLangs{AuthorID: "testAuthor"}},
			[]LangView{{
				ID:   "testId",
				Name: "EN",
			}},
			false,
		},
	}

	for _, tt := range tests {
		v := validator.New()
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewAllLangsHandler(f.langRepo, v)
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
