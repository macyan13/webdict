package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddTranslationHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		translationRepo translation.Repository
		validator       validator
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
			"Error on validation",
			func() fields {
				return fields{
					validator: newFailValidator(),
				}
			},
			args{cmd: AddTranslation{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"Error on translation creation",
			func() fields {
				return fields{
					validator: newSuccessValidator(),
				}
			},
			args{cmd: AddTranslation{Source: "text", Target: "test", AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"Error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				return fields{
					translationRepo: &translationRepo,
					validator:       newSuccessValidator(),
				}
			},
			args{cmd: AddTranslation{Source: "text", Target: "test", AuthorID: "testAuthor", LangID: "testLang"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := AddTranslationHandler{
				translationRepo: tt.fieldsFn().translationRepo,
				validator:       tt.fieldsFn().validator,
			}
			id, err := h.Handle(tt.args.cmd)
			assert.Equal(t, "", id)
			assert.True(t, tt.wantErr(t, err))
		})
	}
}

func TestAddTranslationHandler_Handle_PositiveCase(t *testing.T) {
	tags := []string{"tag1", "tag2"}
	authorID := "testAuthor"
	source := "text"
	langID := "testLang"

	translationRepo := translation.MockRepository{}
	translationRepo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(nil)

	handler := AddTranslationHandler{
		translationRepo: &translationRepo,
		validator:       newSuccessValidator(),
	}

	cmd := AddTranslation{
		Transcription: "transcription",
		Target:        "target",
		Source:        source,
		Example:       "example",
		TagIds:        tags,
		AuthorID:      authorID,
		LangID:        langID,
	}

	id, err := handler.Handle(cmd)
	assert.Nil(t, err)

	createdTranslation := translationRepo.Calls[0].Arguments[0].(*translation.Translation)
	data := createdTranslation.ToMap()

	assert.Equal(t, id, createdTranslation.ID())
	assert.Equal(t, cmd.Target, data["target"])
	assert.Equal(t, cmd.Transcription, data["transcription"])
	assert.Equal(t, cmd.Source, data["source"])
	assert.Equal(t, cmd.TagIds, data["tagIDs"])
	assert.Equal(t, cmd.AuthorID, data["authorID"])
	assert.Equal(t, cmd.LangID, data["langID"])
}
