package query

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
)

func TestAllTagsHandler_Handle(t *testing.T) {
	type fields struct {
		tagRepo TagViewRepository
	}
	type args struct {
		cmd AllTags
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     []TagView
		wantErr  bool
	}{
		{
			"Error on validation",
			func() fields {
				return fields{tagRepo: &MockTagViewRepository{}}
			},
			args{cmd: AllTags{AuthorID: ""}},
			nil,
			true,
		},
		{
			"Error on DB query",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return(nil, errors.New("testErr"))
				return fields{tagRepo: &repo}
			},
			args{cmd: AllTags{AuthorID: "testAuthor"}},
			nil,
			true,
		},
		{
			"Positive case",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return([]TagView{{
					ID:   "testId",
					Name: "test Tag",
				}}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: AllTags{AuthorID: "testAuthor"}},
			[]TagView{{
				ID:   "testId",
				Name: "test Tag",
			}},
			false,
		},
		{
			"Check sanitization",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return([]TagView{{
					ID:   "testId",
					Name: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')"><br>Test Tag</br><a>`,
				}}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: AllTags{AuthorID: "testAuthor"}},
			[]TagView{{
				ID:   "testId",
				Name: "Test Tag",
			}},
			false,
		},
	}

	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewAllTagsHandler(fields.tagRepo, v)
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
