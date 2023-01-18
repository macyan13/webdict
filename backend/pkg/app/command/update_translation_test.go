package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateTranslationHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		translationRepo translation.Repository
		tagRepo         tag.Repository
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
			"Case 1: Can not get translation from DB",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(nil, errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tag.MockRepository{},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Case 2: Tags repo can not perform query",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(&translation.Translation{}, nil)
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tagRepo,
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Case 3: Tags not exist",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(&translation.Translation{}, nil)
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, nil)
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tagRepo,
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not apply changes for translation testID, some passed tag are not found", err.Error(), i)
				return true
			},
		},
		{
			"Case 4: tags no set, error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				tr, err := translation.NewTranslation("new", "new", "new", "new", "new", []string{})
				assert.Nil(t, err)

				translationRepo.On("Get", "testID", "testAuthor").Return(tr, nil)
				translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tag.MockRepository{},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", Text: "test", Translation: "test", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Case 5: error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				tr, err := translation.NewTranslation("new", "new", "new", "new", "new", []string{})
				assert.Nil(t, err)
				translationRepo.On("Get", "testID", "testAuthor").Return(tr, nil)
				translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(true, nil)
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tagRepo,
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", Text: "test", Translation: "test", TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewUpdateTranslationHandler(fields.translationRepo, fields.tagRepo)
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestUpdateTranslationHandler_Handle_PositiveCase(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	authorID := "testAuthor"
	id := "testID"

	tagRepo := tag.MockRepository{}
	tagRepo.On("AllExist", tags, authorID).Return(true, nil)

	translationRepo := translation.MockRepository{}
	tr, err := translation.NewTranslation("test", "", "test", authorID, "", []string{})
	assert.Nil(t, err)
	translationRepo.On("Get", id, authorID).Return(tr, nil)
	translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(nil)

	handler := NewUpdateTranslationHandler(
		&translationRepo,
		&tagRepo,
	)

	cmd := UpdateTranslation{
		ID:            id,
		Transcription: "transcription",
		Translation:   "translation",
		Text:          "text",
		Example:       "example",
		TagIds:        tags,
		AuthorID:      authorID,
	}
	assert.Nil(t, handler.Handle(cmd))

	updatedTranslation := translationRepo.Calls[1].Arguments[0].(*translation.Translation)
	data := updatedTranslation.ToMap()

	assert.Equal(t, cmd.Translation, data["translation"])
	assert.Equal(t, cmd.Transcription, data["transcription"])
	assert.Equal(t, cmd.Text, data["text"])
	assert.Equal(t, cmd.Example, data["example"])
	assert.Equal(t, cmd.TagIds, data["tagIDs"])
}
