package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateTranslationHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		translationRepo translation.Repository
		validator       validator
	}
	type args struct {
		cmd UpdateTranslation
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on validation",
			func() fields {
				return fields{
					validator: newFailValidator(),
				}
			},
			args{cmd: UpdateTranslation{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"Error on getting translation",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(nil, errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					validator:       newSuccessValidator(),
				}
			},
			args{cmd: UpdateTranslation{TagIds: []string{"tag1"}, AuthorID: "testAuthor", ID: "testID"}},
			assert.Error,
		},
		{
			"Error on applying changes",
			func() fields {
				translationRepo := translation.MockRepository{}
				tr, err := translation.NewTranslation("new", "new", "new", "new", "new", []string{}, "new")
				assert.Nil(t, err)
				translationRepo.On("Get", "testID", "testAuthor").Return(tr, nil)
				return fields{
					translationRepo: &translationRepo,
					validator:       newSuccessValidator(),
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", Source: "test", Target: "test", TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"error on update",
			func() fields {
				translationRepo := translation.MockRepository{}
				tr, err := translation.NewTranslation("new", "new", "new", "new", "new", []string{}, "new")
				assert.Nil(t, err)
				translationRepo.On("Get", "testID", "testAuthor").Return(tr, nil)
				translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					validator:       newSuccessValidator(),
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", Source: "test", Target: "test", TagIds: []string{"tag1"}, AuthorID: "testAuthor", LangID: "langID"}},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateTranslationHandler{
				translationRepo: tt.fieldsFn().translationRepo,
				validator:       tt.fieldsFn().validator,
			}
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestUpdateTranslationHandler_Handle_PositiveCase(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	authorID := "testAuthor"
	id := "testID"
	langID := "langID"

	translationRepo := translation.MockRepository{}
	tr, err := translation.NewTranslation("test", "", "test", authorID, "", []string{}, "new")
	assert.Nil(t, err)
	translationRepo.On("Get", id, authorID).Return(tr, nil)
	translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(nil)

	handler := UpdateTranslationHandler{
		translationRepo: &translationRepo,
		validator:       newSuccessValidator(),
	}

	cmd := UpdateTranslation{
		ID:            id,
		Target:        "transcription",
		Transcription: "translation",
		Source:        "text",
		Example:       "example",
		TagIds:        tags,
		AuthorID:      authorID,
		LangID:        langID,
	}
	assert.Nil(t, handler.Handle(cmd))

	updatedTranslation := translationRepo.Calls[1].Arguments[0].(*translation.Translation)
	data := updatedTranslation.ToMap()

	assert.Equal(t, cmd.Transcription, data["transcription"])
	assert.Equal(t, cmd.Target, data["target"])
	assert.Equal(t, cmd.Source, data["source"])
	assert.Equal(t, cmd.Example, data["example"])
	assert.Equal(t, cmd.TagIds, data["tagIDs"])
	assert.Equal(t, cmd.LangID, data["langID"])
}
