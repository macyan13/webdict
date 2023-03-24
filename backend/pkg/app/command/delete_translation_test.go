package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteTranslationHandler_Handle(t *testing.T) {
	type fields struct {
		translationRepo translation.Repository
	}
	type args struct {
		cmd DeleteTranslation
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Case 1: error during translation removing",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Delete", "testID", "testAuthor").Return(errors.New("testErr"))
				return fields{translationRepo: &repo}
			},
			args{cmd: DeleteTranslation{
				ID:       "testID",
				AuthorID: "testAuthor",
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Case 2: positive case",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Delete", "testID", "testAuthor").Return(nil)
				return fields{translationRepo: &repo}
			},
			args{cmd: DeleteTranslation{
				ID:       "testID",
				AuthorID: "testAuthor",
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
			h := NewDeleteTranslationHandler(fields.translationRepo)
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}
