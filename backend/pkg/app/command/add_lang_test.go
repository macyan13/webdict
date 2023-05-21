package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAddLangHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		langRepo lang.Repository
	}
	type args struct {
		cmd AddLang
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on lang creation",
			func() fields {
				return fields{langRepo: &lang.MockRepository{}}
			},
			args{cmd: AddLang{
				Name:     "",
				AuthorID: "testAuthor",
			}},
			assert.Error,
		},
		{
			"Error on lang saving",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("ExistByName", "en", "testAuthor").Return(false, nil)
				langRepo.On("Create", mock.AnythingOfType("*lang.Lang")).Return(errors.New("testError"))
				return fields{langRepo: &langRepo}
			},
			args{cmd: AddLang{
				Name:     "en",
				AuthorID: "testAuthor",
			}},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := NewAddLangHandler(f.langRepo)
			id, err := h.Handle(tt.args.cmd)
			assert.Equal(t, "", id)
			tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestAddLangHandler_Handle_PositiveCase(t *testing.T) {
	ln := "en"
	authorID := "testAuthor"
	langRepo := lang.MockRepository{}
	langRepo.On("Create", mock.AnythingOfType("*lang.Lang")).Return(nil)

	handler := NewAddLangHandler(&langRepo)
	cmd := AddLang{
		Name:     ln,
		AuthorID: authorID,
	}

	id, err := handler.Handle(cmd)
	assert.Nil(t, err)

	createdLang := langRepo.Calls[0].Arguments[0].(*lang.Lang)
	data := createdLang.ToMap()

	assert.Equal(t, createdLang.ID(), id)
	assert.Equal(t, cmd.Name, data["name"])
	assert.Equal(t, cmd.AuthorID, data["authorID"])
}
