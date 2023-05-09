package command

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func Test_validator_validateTags(t *testing.T) {
	type fields struct {
		tagRepo tag.Repository
	}
	type args struct {
		data translationData
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
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, fmt.Errorf("testErr"))
				return fields{
					tagRepo: &tagRepo,
				}
			},
			args{data: translationData{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"Tags not exist",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(false, nil)
				return fields{
					tagRepo: &tagRepo,
				}
			},
			args{data: translationData{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "some of passed tags: [tag1] are not found", err.Error(), i)
				return true
			},
		},
		{
			"Tags exist",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tag1"}, "testAuthor").Return(true, nil)
				return fields{
					tagRepo: &tagRepo,
				}
			},
			args{data: translationData{TagIds: []string{"tag1"}, AuthorID: "testAuthor"}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator{
				tagRepo: tt.fieldsFn().tagRepo,
			}
			tt.wantErr(t, v.validateTags(tt.args.data), fmt.Sprintf("validateTags(%v)", tt.args.data))
		})
	}
}

func Test_validator_validateLang(t *testing.T) {
	type fields struct {
		langRepo lang.Repository
	}
	type args struct {
		data translationData
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Lang repo can not perform query",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "testAuthor").Return(false, fmt.Errorf("testErr"))
				return fields{
					langRepo: &langRepo,
				}
			},
			args{data: translationData{LangID: "langID", AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"Lang not exist",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "testAuthor").Return(false, nil)
				return fields{
					langRepo: &langRepo,
				}
			},
			args{data: translationData{LangID: "langID", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "lang with id: langID is not found", err.Error(), i)
				return true
			},
		},
		{
			"Lang exist",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "testAuthor").Return(true, nil)
				return fields{
					langRepo: &langRepo,
				}
			},
			args{data: translationData{LangID: "langID", AuthorID: "testAuthor"}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator{
				langRepo: tt.fieldsFn().langRepo,
			}
			tt.wantErr(t, v.validateLang(tt.args.data), fmt.Sprintf("validateLang(%v)", tt.args.data))
		})
	}
}

func Test_validator_validateSource(t *testing.T) {
	type fields struct {
		translationRepo translation.Repository
	}
	type args struct {
		data translationData
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Translation repo can not perform query",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "source", "testAuthor", "langID").Return(false, fmt.Errorf("testErr"))
				return fields{
					translationRepo: &translationRepo,
				}
			},
			args{data: translationData{Source: "source", LangID: "langID", AuthorID: "testAuthor"}},
			assert.Error,
		},
		{
			"Translation exist",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "source", "testAuthor", "langID").Return(true, nil)
				return fields{
					translationRepo: &translationRepo,
				}
			},
			args{data: translationData{Source: "source", LangID: "langID", AuthorID: "testAuthor"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "translation with source: source already created", err.Error(), i)
				return true
			},
		},
		{
			"Translation not exist",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "source", "testAuthor", "langID").Return(false, nil)
				return fields{
					translationRepo: &translationRepo,
				}
			},
			args{data: translationData{Source: "source", LangID: "langID", AuthorID: "testAuthor"}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator{
				translationRepo: tt.fieldsFn().translationRepo,
			}
			tt.wantErr(t, v.validateSource(tt.args.data), fmt.Sprintf("validateSource(%v)", tt.args.data))
		})
	}
}

func Test_validator_validate(t *testing.T) {
	type fields struct {
		tagRepo         tag.Repository
		langRepo        lang.Repository
		translationRepo translation.Repository
	}
	type args struct {
		data translationData
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Validate multiple errors",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tagID"}, "testAuthor").Return(false, nil)
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "testAuthor").Return(false, nil)
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "source", "testAuthor", "langID").Return(true, nil)
				return fields{
					tagRepo:         &tagRepo,
					langRepo:        &langRepo,
					translationRepo: &translationRepo,
				}
			},
			args{data: translationData{TagIds: []string{"tagID"}, LangID: "langID", AuthorID: "testAuthor", Source: "source"}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "some of passed tags: [tagID] are not found"), i)
				assert.True(t, strings.Contains(err.Error(), "lang with id: langID is not found"), i)
				assert.True(t, strings.Contains(err.Error(), "translation with source: source already created"), i)
				return true
			},
		},
		{
			"Validate no errors",
			func() fields {
				tagRepo := tag.MockRepository{}
				tagRepo.On("AllExist", []string{"tagID"}, "testAuthor").Return(true, nil)
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "testAuthor").Return(true, nil)
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistBySource", "source", "testAuthor", "langID").Return(false, nil)
				return fields{
					tagRepo:         &tagRepo,
					langRepo:        &langRepo,
					translationRepo: &translationRepo,
				}
			},
			args{data: translationData{TagIds: []string{"tagID"}, LangID: "langID", AuthorID: "testAuthor", Source: "source"}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator{
				tagRepo:         tt.fieldsFn().tagRepo,
				langRepo:        tt.fieldsFn().langRepo,
				translationRepo: tt.fieldsFn().translationRepo,
			}
			tt.wantErr(t, v.validate(tt.args.data), fmt.Sprintf("validate(%v)", tt.args.data))
		})
	}
}

func newFailValidator() validator {
	tagRepo := tag.MockRepository{}
	tagRepo.On("AllExist", mock.Anything, mock.Anything).Return(false, fmt.Errorf("testErr"))
	langRepo := lang.MockRepository{}
	langRepo.On("Exist", mock.Anything, mock.Anything).Return(false, fmt.Errorf("testErr"))
	translationRepo := translation.MockRepository{}
	translationRepo.On("ExistBySource", mock.Anything, mock.Anything, mock.Anything).Return(true, fmt.Errorf("testErr"))
	return validator{
		tagRepo:         &tagRepo,
		langRepo:        &langRepo,
		translationRepo: &translationRepo,
	}
}

func newSuccessValidator() validator {
	tagRepo := tag.MockRepository{}
	tagRepo.On("AllExist", mock.Anything, mock.Anything).Return(true, nil)
	langRepo := lang.MockRepository{}
	langRepo.On("Exist", mock.Anything, mock.Anything).Return(true, nil)
	translationRepo := translation.MockRepository{}
	translationRepo.On("ExistBySource", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	return validator{
		tagRepo:         &tagRepo,
		langRepo:        &langRepo,
		translationRepo: &translationRepo,
	}
}
