package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteUserHandler_Handle(t *testing.T) {
	type fields struct {
		userRepo        user.Repository
		langRepo        lang.Repository
		tagRepo         tag.Repository
		translationRepo translation.Repository
	}
	type args struct {
		cmd DeleteUser
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     int
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on user delete",
			func() fields {
				userRepo := user.NewMockRepository(t)
				userRepo.On("Delete", "authorID").Return(0, errors.New("test"))
				tagRepo := tag.NewMockRepository(t)
				tagRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				langRepo := lang.NewMockRepository(t)
				langRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				translationRepo := translation.NewMockRepository(t)
				translationRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				return fields{
					userRepo:        userRepo,
					langRepo:        langRepo,
					tagRepo:         tagRepo,
					translationRepo: translationRepo,
				}
			},
			args{cmd: DeleteUser{AuthorID: "authorID"}},
			3,
			assert.Error,
		},
		{
			"Error on tag delete",
			func() fields {
				userRepo := user.NewMockRepository(t)
				userRepo.On("Delete", "authorID").Return(1, nil)
				tagRepo := tag.NewMockRepository(t)
				tagRepo.On("DeleteByAuthorID", "authorID").Return(0, errors.New("test"))
				langRepo := lang.NewMockRepository(t)
				langRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				translationRepo := translation.NewMockRepository(t)
				translationRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				return fields{
					userRepo:        userRepo,
					langRepo:        langRepo,
					tagRepo:         tagRepo,
					translationRepo: translationRepo,
				}
			},
			args{cmd: DeleteUser{AuthorID: "authorID"}},
			3,
			assert.Error,
		},
		{
			"Error on lang delete",
			func() fields {
				userRepo := user.NewMockRepository(t)
				userRepo.On("Delete", "authorID").Return(1, nil)
				tagRepo := tag.NewMockRepository(t)
				tagRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				langRepo := lang.NewMockRepository(t)
				langRepo.On("DeleteByAuthorID", "authorID").Return(0, errors.New("test"))
				translationRepo := translation.NewMockRepository(t)
				translationRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				return fields{
					userRepo:        userRepo,
					langRepo:        langRepo,
					tagRepo:         tagRepo,
					translationRepo: translationRepo,
				}
			},
			args{cmd: DeleteUser{AuthorID: "authorID"}},
			3,
			assert.Error,
		},
		{
			"Error on tag translation",
			func() fields {
				userRepo := user.NewMockRepository(t)
				userRepo.On("Delete", "authorID").Return(1, nil)
				tagRepo := tag.NewMockRepository(t)
				tagRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				langRepo := lang.NewMockRepository(t)
				langRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				translationRepo := translation.NewMockRepository(t)
				translationRepo.On("DeleteByAuthorID", "authorID").Return(0, errors.New("test"))
				return fields{
					userRepo:        userRepo,
					langRepo:        langRepo,
					tagRepo:         tagRepo,
					translationRepo: translationRepo,
				}
			},
			args{cmd: DeleteUser{AuthorID: "authorID"}},
			3,
			assert.Error,
		},
		{
			"Everything removed without errors",
			func() fields {
				userRepo := user.NewMockRepository(t)
				userRepo.On("Delete", "authorID").Return(1, nil)
				tagRepo := tag.NewMockRepository(t)
				tagRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				langRepo := lang.NewMockRepository(t)
				langRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				translationRepo := translation.NewMockRepository(t)
				translationRepo.On("DeleteByAuthorID", "authorID").Return(1, nil)
				return fields{
					userRepo:        userRepo,
					langRepo:        langRepo,
					tagRepo:         tagRepo,
					translationRepo: translationRepo,
				}
			},
			args{cmd: DeleteUser{AuthorID: "authorID"}},
			4,
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			h := &DeleteUserHandler{
				userRepo:        f.userRepo,
				langRepo:        f.langRepo,
				tagRepo:         f.tagRepo,
				translationRepo: f.translationRepo,
			}
			got, err := h.Handle(tt.args.cmd)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.cmd)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v)", tt.args.cmd)
		})
	}
}
