package command

import (
	"errors"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestAddTranslationHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		translationRepo    translation.Repository
		tagRepo            tag.Repository
		supportedLanguages []translation.Lang
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
				assert.Equal(t, "some passed tag are not found", err.Error(), i)
				return true
			},
		},
		{
			"Target repo can not perform query",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "text", "testAuthor").Return(false, errors.New("testErr"))
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
			"Target already exists",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "text", "testAuthor").Return(true, nil)
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
			"Lang is not configured",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "text", "testAuthor").Return(false, nil)
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tag.MockRepository{},
					supportedLanguages: []translation.Lang{"DE"},
				}
			},
			args{cmd: AddTranslation{Text: "text", AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "passed language EN is not supported", err.Error(), i)
				return true
			},
		},
		{
			"Error on Apply changes",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "text", "testAuthor").Return(false, nil)
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tag.MockRepository{},
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: AddTranslation{Text: "text", AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "target can not be empty"), i)
				return true
			},
		},
		{
			"Error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "text", "testAuthor").Return(false, nil)
				translationRepo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tag.MockRepository{},
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: AddTranslation{Text: "text", Target: "test", AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewAddTranslationHandler(
				f.translationRepo,
				f.tagRepo,
				f.supportedLanguages,
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
	source := "text"

	tagRepo := tag.MockRepository{}
	tagRepo.On("AllExist", tags, authorID).Return(true, nil)

	translationRepo := translation.MockRepository{}
	translationRepo.On("ExistBySource", source, authorID).Return(false, nil)
	translationRepo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(nil)

	handler := NewAddTranslationHandler(
		&translationRepo,
		&tagRepo,
		[]translation.Lang{"EN"},
	)

	cmd := AddTranslation{
		Transcription: "transcription",
		Target:        "target",
		Text:          source,
		Example:       "example",
		TagIds:        tags,
		AuthorID:      "testAuthor",
		Lang:          translation.Lang("EN"),
	}

	id, err := handler.Handle(cmd)
	assert.Nil(t, err)

	createdTranslation := translationRepo.Calls[1].Arguments[0].(*translation.Translation)
	data := createdTranslation.ToMap()

	assert.Equal(t, id, createdTranslation.ID())
	assert.Equal(t, cmd.Target, data["target"])
	assert.Equal(t, cmd.Transcription, data["transcription"])
	assert.Equal(t, cmd.Text, data["source"])
	assert.Equal(t, cmd.TagIds, data["tagIDs"])
	assert.Equal(t, cmd.AuthorID, data["authorID"])
	assert.Equal(t, string(cmd.Lang), data["lang"])
}
