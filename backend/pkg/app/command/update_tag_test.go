package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateTagHandler_Handle(t *testing.T) {
	type fields struct {
		tagRepo tag.Repository
	}
	type args struct {
		cmd UpdateTag
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Case 1: can not get tag from DB",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("Get", "testID", "testAuthor").Return(nil, errors.New("testError"))
				return fields{tagRepo: &tagRepo}
			},
			args{cmd: UpdateTag{
				TagID:    "testID",
				Tag:      "tag",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return true
			},
		},
		{
			"Case 2: error on saving",
			func() fields {
				tg := tag.UnmarshalFromDB("testID", "testTag", "testAuthor")
				tagRepo := tag.MockRepository{}
				tagRepo.On("Get", "testID", "testAuthor").Return(tg, nil)

				updatedTg := tag.UnmarshalFromDB("testID", "updatedTag", "testAuthor")
				tagRepo.On("Update", updatedTg).Return(errors.New("testError"))
				return fields{tagRepo: &tagRepo}
			},
			args{cmd: UpdateTag{
				TagID:    "testID",
				Tag:      "updatedTag",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return true
			},
		},
		{
			"Case 3: positive case",
			func() fields {
				tg := tag.UnmarshalFromDB("testID", "testTag", "testAuthor")
				tagRepo := tag.MockRepository{}
				tagRepo.On("Get", "testID", "testAuthor").Return(tg, nil)

				updatedTg := tag.UnmarshalFromDB("testID", "updatedTag", "testAuthor")
				tagRepo.On("Update", updatedTg).Return(nil)
				return fields{tagRepo: &tagRepo}
			},
			args{cmd: UpdateTag{
				TagID:    "testID",
				Tag:      "updatedTag",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUpdateTagHandler(tt.fieldsFn().tagRepo)
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}
