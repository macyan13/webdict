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

	err = usr.ApplyChanges(name, email, password, role, usr.DefaultLangID(), user.NewListOptions(true))
	assert.Nil(t, err)

	repo := UserRepo{}

	model, err := repo.fromDomainToModel(usr)
	assert.Nil(t, err)
	assert.Equal(t, usr.ID(), model.ID)
	assert.Equal(t, name, name)
	assert.Equal(t, email, model.Email)
	assert.Equal(t, password, model.Password)
	assert.Equal(t, int(role), model.Role)
	assert.Equal(t, true, model.ListOptions.ShowTranscription)
}

func TestUserRepo_fromModelToView(t *testing.T) {
	type fields struct {
		langRepo      query.LangViewRepository
		roleConverter *query.RoleConverter
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
				return fields{langRepo: &query.MockLangViewRepository{}, roleConverter: query.NewRoleMapper()}
			},
			args{model: UserModel{ID: "authorID", Role: 1}},
			query.UserView{ID: "authorID", Role: query.RoleView{
				ID:      1,
				Name:    "Admin",
				IsAdmin: true,
			}},
			assert.NoError,
		},
		{
			"Error on getting lang view",
			func() fields {
				langRepo := query.MockLangViewRepository{}
				langRepo.On("GetView", "langID", "authorID").Return(query.LangView{}, fmt.Errorf("testError"))
				return fields{langRepo: &langRepo, roleConverter: query.NewRoleMapper()}
			},
			args{model: UserModel{ID: "authorID", DefaultLangID: "langID", Role: 1}},
			query.UserView{},
			assert.Error,
		},
		{
			"Error on converting role view",
			func() fields {
				return fields{langRepo: &query.MockLangViewRepository{}, roleConverter: query.NewRoleMapper()}
			},
			args{model: UserModel{ID: "authorID", DefaultLangID: "langID", Role: 3}},
			query.UserView{},
			assert.Error,
		},
		{
			"Positive case",
			func() fields {
				langRepo := query.MockLangViewRepository{}
				langRepo.On("GetView", "langID", "authorID").Return(query.LangView{Name: "test"}, nil)
				return fields{langRepo: &langRepo, roleConverter: query.NewRoleMapper()}
			},
			args{model: UserModel{ID: "authorID", DefaultLangID: "langID", Role: 1, ListOptions: ListOptionsModel{ShowTranscription: true}}},
			query.UserView{ID: "authorID", DefaultLang: query.LangView{Name: "test"}, Role: query.RoleView{
				ID:      1,
				Name:    "Admin",
				IsAdmin: true,
			}, ListOptions: query.UserListOptionsView{ShowTranscription: true}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			r := &UserRepo{
				langViewRepo:  f.langRepo,
				roleConverter: f.roleConverter,
			}
			got, err := r.fromModelToView(tt.args.model)
			if !tt.wantErr(t, err, fmt.Sprintf("fromModelToView(%v)", tt.args.model)) {
				return
			}
			assert.Equalf(t, tt.want, got, "fromModelToView(%v)", tt.args.model)
		})
	}
}

func TestUserRepo_fromModelToDomain(t *testing.T) {
	model := UserModel{
		ID:       "authorID",
		Name:     "John",
		Email:    "John@do.com",
		Password: "testPassword",
		Role:     1,
	}

	repo := UserRepo{}
	usr := repo.fromModelToDomain(model)
	listOptions := usr.ListOptions()

	assert.Equal(t, model.ID, usr.ID())
	assert.Equal(t, model.Email, usr.Email())
	assert.Equal(t, model.Password, usr.Password())
	assert.Equal(t, user.Role(model.Role), usr.Role())
	assert.Equal(t, model.DefaultLangID, usr.DefaultLangID())
	assert.Equal(t, false, listOptions.ToMap()["showTranscription"])
}
