package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddTagHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		tagRepo tag.Repository
	}
	type args struct {
		cmd AddTag
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on tag creation",
			func() fields {
				return fields{tagRepo: &tag.MockRepository{}}
			},
			args{cmd: AddTag{
				Tag:      "t",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "tag length should be at least 2 symbols, 1 passed (t)", err.Error(), i)
				return true
			},
		},
		{
			"Error on checking existing tag",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("ExistByTag", "testTag", "testAuthor").Return(false, errors.New("testError"))
				return fields{tagRepo: &tagRepo}
			},
			args{cmd: AddTag{
				Tag:      "testTag",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return true
			},
		},
		{
			"Tag already exists",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("ExistByTag", "testTag", "testAuthor").Return(true, nil)
				return fields{tagRepo: &tagRepo}
			},
			args{cmd: AddTag{
				Tag:      "testTag",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not create new tag - tag:testTag already created", err.Error(), i)
				return true
			},
		},
		{
			"Error on tag saving",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("ExistByTag", "testTag", "testAuthor").Return(false, nil)
				tagRepo.On("Create", mock.AnythingOfType("*tag.Tag")).Return(errors.New("testError"))
				return fields{tagRepo: &tagRepo}
			},
			args{cmd: AddTag{
				Tag:      "testTag",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewAddTagHandler(fields.tagRepo)
			id, err := h.Handle(tt.args.cmd)
			assert.Equal(t, "", id)
			tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestAddTagHandler_Handle_PositiveCase(t *testing.T) {
	tg := "testTag"
	authorID := "testAuthor"
	tagRepo := tag.MockRepository{}
	tagRepo.On("ExistByTag", tg, authorID).Return(false, nil)
	tagRepo.On("Create", mock.AnythingOfType("*tag.Tag")).Return(nil)

	handler := NewAddTagHandler(&tagRepo)
	cmd := AddTag{
		Tag:      tg,
		AuthorID: authorID,
	}

	id, err := handler.Handle(cmd)
	assert.Nil(t, err)

	createdTag := tagRepo.Calls[1].Arguments[0].(*tag.Tag)
	data := createdTag.ToMap()

	assert.Equal(t, createdTag.ID(), id)
	assert.Equal(t, cmd.Tag, data["tag"])
	assert.Equal(t, cmd.AuthorID, data["authorID"])
}
