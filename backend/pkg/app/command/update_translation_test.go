package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateTranslationHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		translationRepo    translation.Repository
		tagRepo            tag.Repository
		supportedLanguages []translation.Lang
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
			"Can not get translation from DB",
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
			"Tags repo can not perform query",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(&translation.Translation{}, nil)
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, errors.New("testErr"))
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tagRepo,
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", TagIds: []string{"tag1"}, AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Lang is not supported",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(&translation.Translation{}, nil)
				return fields{
					translationRepo:    &translationRepo,
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", TagIds: []string{"tag1"}, AuthorID: "testAuthor", Lang: translation.Lang("DE")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "passed language DE is not supported", err.Error(), i)
				return true
			},
		},
		{
			"Tags not exist",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("Get", "testID", "testAuthor").Return(&translation.Translation{}, nil)
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, nil)
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tagRepo,
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", TagIds: []string{"tag1"}, AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not apply changes for translation testID, some passed tag are not found", err.Error(), i)
				return true
			},
		},
		{
			"tags no set, error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				tr, err := translation.NewTranslation("new", "new", "new", "new", "new", []string{}, translation.Lang("EN"))
				assert.Nil(t, err)

				translationRepo.On("Get", "testID", "testAuthor").Return(tr, nil)
				translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tag.MockRepository{},
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", Source: "test", Target: "test", AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"error on save",
			func() fields {
				translationRepo := translation.MockRepository{}
				tr, err := translation.NewTranslation("new", "new", "new", "new", "new", []string{}, translation.Lang("EN"))
				assert.Nil(t, err)
				translationRepo.On("Get", "testID", "testAuthor").Return(tr, nil)
				translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(errors.New("testErr"))
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(true, nil)
				return fields{
					translationRepo:    &translationRepo,
					tagRepo:            &tagRepo,
					supportedLanguages: []translation.Lang{"EN"},
				}
			},
			args{cmd: UpdateTranslation{ID: "testID", Source: "test", Target: "test", TagIds: []string{"tag1"}, AuthorID: "testAuthor", Lang: translation.Lang("EN")}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewUpdateTranslationHandler(f.translationRepo, f.tagRepo, f.supportedLanguages)
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
	tr, err := translation.NewTranslation("test", "", "test", authorID, "", []string{}, translation.Lang("EN"))
	assert.Nil(t, err)
	translationRepo.On("Get", id, authorID).Return(tr, nil)
	translationRepo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(nil)

	handler := NewUpdateTranslationHandler(
		&translationRepo,
		&tagRepo,
		[]translation.Lang{"EN"},
	)

	cmd := UpdateTranslation{
		ID:            id,
		Target:        "transcription",
		Transcription: "translation",
		Source:        "text",
		Example:       "example",
		TagIds:        tags,
		AuthorID:      authorID,
		Lang:          translation.Lang("EN"),
	}
	assert.Nil(t, handler.Handle(cmd))

	updatedTranslation := translationRepo.Calls[1].Arguments[0].(*translation.Translation)
	data := updatedTranslation.ToMap()

	assert.Equal(t, cmd.Transcription, data["transcription"])
	assert.Equal(t, cmd.Target, data["target"])
	assert.Equal(t, cmd.Source, data["source"])
	assert.Equal(t, cmd.Example, data["example"])
	assert.Equal(t, cmd.TagIds, data["tagIDs"])
	assert.Equal(t, string(cmd.Lang), data["lang"])
}
