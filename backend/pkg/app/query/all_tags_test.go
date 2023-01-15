package query

import (
	"errors"
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
			"Case 1: error on DB query",
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
			"Case 2: positive case",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return([]TagView{{
					ID:  "testId",
					Tag: "test Tag",
				}}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: AllTags{AuthorID: "testAuthor"}},
			[]TagView{{
				ID:  "testId",
				Tag: "test Tag",
			}},
			false,
		},
		{
			"Case 3: check sanitization",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetAllViews", "testAuthor").Return([]TagView{{
					ID:  "testId",
					Tag: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">Test Tag<a>`,
				}}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: AllTags{AuthorID: "testAuthor"}},
			[]TagView{{
				ID:  "testId",
				Tag: "Test Tag",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewAllTagsHandler(fields.tagRepo)
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
