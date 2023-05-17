package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteLangHandler_Handle(t *testing.T) {
	type fields struct {
		langRepo        lang.Repository
		translationRepo translation.Repository
	}
	type args struct {
		cmd DeleteLang
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Translation repo returns error on validation",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByLang", "testId", "testAuthorID").Return(false, errors.New("testError"))
				return fields{
					langRepo:        &lang.MockRepository{},
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteLang{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			assert.Error,
		},
		{
			"Translation with the lang exist",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByLang", "testId", "testAuthorID").Return(true, nil)
				return fields{
					langRepo:        &lang.MockRepository{},
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteLang{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not remove lang:testId as some translations use it", err.Error(), i)
				return true
			},
		},
		{
			"LangID repo returns error",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByLang", "testId", "testAuthorID").Return(false, nil)
				langRepo := lang.MockRepository{}
				langRepo.On("Delete", "testId", "testAuthorID").Return(errors.New("testError"))
				return fields{
					langRepo:        &langRepo,
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteLang{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			assert.Error,
		},
		{
			"Positive",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByLang", "testId", "testAuthorID").Return(false, nil)
				langRepo := lang.MockRepository{}
				langRepo.On("Delete", "testId", "testAuthorID").Return(nil)
				return fields{
					langRepo:        &langRepo,
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteLang{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewDeleteLangHandler(
				f.langRepo,
				f.translationRepo,
			)
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}
