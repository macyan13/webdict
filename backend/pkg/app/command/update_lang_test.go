package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateLangHandler_Handle(t *testing.T) {
	type fields struct {
		langRepo lang.Repository
	}
	type args struct {
		cmd UpdateLang
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Can not get lang from DB",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Get", "testID", "testAuthor").Return(nil, errors.New("testError"))
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateLang{
				ID:       "testID",
				Name:     "en",
				AuthorID: "testAuthor",
			}},
			assert.Error,
		},
		{
			"Error on saving",
			func() fields {
				ln := lang.UnmarshalFromDB("testID", "en", "testAuthor")
				langRepo := lang.MockRepository{}
				langRepo.On("Get", "testID", "testAuthor").Return(ln, nil)

				updatedLn := lang.UnmarshalFromDB("testID", "de", "testAuthor")
				langRepo.On("Update", updatedLn).Return(errors.New("testError"))
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateLang{
				ID:       "testID",
				Name:     "de",
				AuthorID: "testAuthor",
			}},
			assert.Error,
		},
		{
			"Error on applying changes",
			func() fields {
				ln := lang.UnmarshalFromDB("testID", "en", "testAuthor")
				langRepo := lang.MockRepository{}
				langRepo.On("Get", "testID", "testAuthor").Return(ln, nil)
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateLang{
				ID:       "testID",
				Name:     "",
				AuthorID: "testAuthor",
			}},
			assert.Error,
		},
		{
			"Positive case",
			func() fields {
				ln := lang.UnmarshalFromDB("testID", "en", "testAuthor")
				langRepo := lang.MockRepository{}
				langRepo.On("Get", "testID", "testAuthor").Return(ln, nil)

				updatedTg := lang.UnmarshalFromDB("testID", "de", "testAuthor")
				langRepo.On("Update", updatedTg).Return(nil)
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateLang{
				ID:       "testID",
				Name:     "de",
				AuthorID: "testAuthor",
			}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUpdateLangHandler(tt.fieldsFn().langRepo)
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}
