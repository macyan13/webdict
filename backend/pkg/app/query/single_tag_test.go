package query

import (
	"errors"
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
			"Case 1: error on DB query",
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
			"Case 2: positive case",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetView", "tagID", "testAuthor").Return(TagView{
					ID:  "tagID",
					Tag: "tag",
				}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: SingleTag{ID: "tagID", AuthorID: "testAuthor"}},
			TagView{
				ID:  "tagID",
				Tag: "tag",
			},
			false,
		},
		{
			"Case 3: check sanitization",
			func() fields {
				repo := MockTagViewRepository{}
				repo.On("GetView", "tagID", "testAuthor").Return(TagView{
					ID:  "tagID",
					Tag: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">Test Tag<a>`,
				}, nil)
				return fields{tagRepo: &repo}
			},
			args{cmd: SingleTag{ID: "tagID", AuthorID: "testAuthor"}},
			TagView{
				ID:  "tagID",
				Tag: "Test Tag",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSingleTagHandler(tt.fieldsFn().tagRepo)
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
