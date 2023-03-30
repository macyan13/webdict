package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteTagHandler_Handle(t *testing.T) {
	type fields struct {
		tagRepo         tag.Repository
		translationRepo translation.Repository
	}
	type args struct {
		cmd DeleteTag
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Case 1: translation repo returns error on validation",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByTag", "testId", "testAuthorID").Return(false, errors.New("testError"))
				return fields{
					tagRepo:         &tag.MockRepository{},
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteTag{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return true
			},
		},
		{
			"Case 2: translation with the tag exist",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByTag", "testId", "testAuthorID").Return(true, nil)
				return fields{
					tagRepo:         &tag.MockRepository{},
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteTag{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not remove tag:testId as some translation is tagged by it", err.Error(), i)
				return true
			},
		},
		{
			"Case 3: tag repo returns error",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByTag", "testId", "testAuthorID").Return(false, nil)
				tagRepo := tag.MockRepository{}
				tagRepo.On("Delete", "testId", "testAuthorID").Return(errors.New("testError"))
				return fields{
					tagRepo:         &tagRepo,
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteTag{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return true
			},
		},
		{
			"Case 4: Positive",
			func() fields {
				translationRepo := translation.MockRepository{}
				translationRepo.On("ExistByTag", "testId", "testAuthorID").Return(false, nil)
				tagRepo := tag.MockRepository{}
				tagRepo.On("Delete", "testId", "testAuthorID").Return(nil)
				return fields{
					tagRepo:         &tagRepo,
					translationRepo: &translationRepo,
				}
			},
			args{cmd: DeleteTag{
				ID:       "testId",
				AuthorID: "testAuthorID",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewDeleteTagHandler(
				fields.tagRepo,
				fields.translationRepo,
			)
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}
