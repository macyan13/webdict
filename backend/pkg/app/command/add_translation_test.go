package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestAddTranslationHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		translationRepo translation.Repository
		tagRepo         tag.Repository
	}
	type args struct {
		cmd AddTranslation
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Tags repo can not perform query",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, errors.New("testErr"))
				return fields{
					translationRepo: &translation.MockRepository{},
					tagRepo:         &tagRepo,
				}
			},
			args{cmd: AddTranslation{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Tags not exist",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, nil)
				return fields{
					translationRepo: &translation.MockRepository{},
					tagRepo:         &tagRepo,
				}
			},
			args{cmd: AddTranslation{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not apply changes for translation tags, some passed tag are not found", err.Error(), i)
				return true
			},
		},
		{
			"Translation repo can not perform query",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByText", "text", "testAuthor").Return(false, errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tag.MockRepository{},
				}
			},
			args{cmd: AddTranslation{Text: "text", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Translation already exists",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByText", "text", "testAuthor").Return(true, nil)
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tag.MockRepository{},
				}
			},
			args{cmd: AddTranslation{Text: "text", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "translation with text: text already created", err.Error(), i)
				return true
			},
		},
		{
			"Error on Apply changes",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByText", "text", "testAuthor").Return(false, nil)
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tag.MockRepository{},
				}
			},
			args{cmd: AddTranslation{Text: "text", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "translation can not be empty"), i)
				return true
			},
		},
		{
			"Error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByText", "text", "testAuthor").Return(false, nil)
				translationRepo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					tagRepo:         &tag.MockRepository{},
				}
			},
			args{cmd: AddTranslation{Text: "text", Translation: "test", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewAddTranslationHandler(
				fields.translationRepo,
				fields.tagRepo,
			)
			id, err := h.Handle(tt.args.cmd)
			assert.Equal(t, "", id)
			assert.True(t, tt.wantErr(t, err))
		})
	}
}

func TestAddTranslationHandler_Handle_PositiveCase(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	authorID := "testAuthor"
	text := "text"

	tagRepo := tag.MockRepository{}
	tagRepo.On("AllExist", tags, authorID).Return(true, nil)

	translationRepo := translation.MockRepository{}
	translationRepo.On("ExistByText", text, authorID).Return(false, nil)
	translationRepo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(nil)

	handler := NewAddTranslationHandler(
		&translationRepo,
		&tagRepo,
	)

	cmd := AddTranslation{
		Transcription: "transcription",
		Translation:   "translation",
		Text:          text,
		Example:       "example",
		TagIds:        tags,
		AuthorID:      "testAuthor",
	}

	id, err := handler.Handle(cmd)
	assert.Nil(t, err)

	createdTranslation := translationRepo.Calls[1].Arguments[0].(*translation.Translation)
	data := createdTranslation.ToMap()

	assert.Equal(t, id, createdTranslation.ID())
	assert.Equal(t, cmd.Translation, data["translation"])
	assert.Equal(t, cmd.Transcription, data["transcription"])
	assert.Equal(t, cmd.Text, data["text"])
	assert.Equal(t, cmd.TagIds, data["tagIDs"])
	assert.Equal(t, cmd.AuthorID, data["authorID"])
	assert.Equal(t, translation.EN, data["lang"])
}
