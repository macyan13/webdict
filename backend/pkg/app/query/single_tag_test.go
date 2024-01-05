package query

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSingleTagHandler_Handle(t *testing.T) {
	type fields struct {
		tagRepo TagViewRepository
	}
	type args struct {
		cmd SingleTag
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     TagView
		wantErr  bool
	}{
		{
			"Error on query validation",
			func() fields {
				return fields{tagRepo: &MockTagViewRepository{}}
			},
			args{cmd: SingleTag{AuthorID: "testAuthor"}},
			TagView{},
			true,
		},
		{
			"Error on DB query",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetView", "tagID", "testAuthor").Return(TagView{}, errors.New("testErr"))
				return fields{tagRepo: &repo}
			},
			args{cmd: SingleTag{ID: "tagID", AuthorID: "testAuthor"}},
			TagView{},
			true,
		},
		{
			"Positive case",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetView", "tagID", "testAuthor").Return(TagView{
					ID:   "tagID",
					Name: "tag",
				}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: SingleTag{ID: "tagID", AuthorID: "testAuthor"}},
			TagView{
				ID:   "tagID",
				Name: "tag",
			},
			false,
		},
		{
			"Check sanitization",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetView", "tagID", "testAuthor").Return(TagView{
					ID:   "tagID",
					Name: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">Test Tag<a>`,
				}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: SingleTag{ID: "tagID", AuthorID: "testAuthor"}},
			TagView{
				ID:   "tagID",
				Name: "Test Tag",
			},
			false,
		},
	}

	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSingleTagHandler(tt.fieldsFn().tagRepo, v)
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

func TestSingleTagValidation(t *testing.T) {
	type args struct {
		query SingleTag
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid case",
			args: args{
				query: SingleTag{
					ID:       "123",
					AuthorID: "456",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Missing ID",
			args: args{
				query: SingleTag{
					AuthorID: "456",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Missing AuthorID",
			args: args{
				query: SingleTag{
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
