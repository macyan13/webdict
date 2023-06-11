package mongo

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepo_fromDomainToModel(t *testing.T) {
	name := "John"
	email := "John@do.com"
	password := "12345678"
	role := user.Admin
	usr, err := user.NewUser(name, email, password, role)

	assert.Nil(t, err)
	repo := UserRepo{}

	model, err := repo.fromDomainToModel(usr)
	assert.Nil(t, err)
	assert.Equal(t, usr.ID(), model.ID)
	assert.Equal(t, name, name)
	assert.Equal(t, email, model.Email)
	assert.Equal(t, password, model.Password)
	assert.Equal(t, int(role), model.Role)
}

func TestUserRepo_fromModelToView(t *testing.T) {
	type fields struct {
		langRepo query.LangViewRepository
	}
	type args struct {
		model UserModel
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     query.UserView
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Default lang is not set",
			func() fields {
				return fields{langRepo: &query.MockLangViewRepository{}}
			},
			args{model: UserModel{ID: "authorID"}},
			query.UserView{ID: "authorID"},
			assert.NoError,
		},
		{
			"Error on getting lang view",
			func() fields {
				langRepo := query.MockLangViewRepository{}
				langRepo.On("GetView", "langID", "authorID").Return(query.LangView{}, fmt.Errorf("testError"))
				return fields{langRepo: &langRepo}
			},
			args{model: UserModel{ID: "authorID", DefaultLangID: "langID"}},
			query.UserView{},
			assert.Error,
		},
		{
			"Positive case",
			func() fields {
				langRepo := query.MockLangViewRepository{}
				langRepo.On("GetView", "langID", "authorID").Return(query.LangView{Name: "test"}, nil)
				return fields{langRepo: &langRepo}
			},
			args{model: UserModel{ID: "authorID", DefaultLangID: "langID"}},
			query.UserView{ID: "authorID", DefaultLang: query.LangView{Name: "test"}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UserRepo{
				langRepo: tt.fieldsFn().langRepo,
			}
			got, err := r.fromModelToView(tt.args.model)
			if !tt.wantErr(t, err, fmt.Sprintf("fromModelToView(%v)", tt.args.model)) {
				return
			}
			assert.Equalf(t, tt.want, got, "fromModelToView(%v)", tt.args.model)
		})
	}
}
