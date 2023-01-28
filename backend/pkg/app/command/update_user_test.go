package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestUpdateUserHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		cipher   Cipher
	}
	type args struct {
		cmd UpdateUser
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Invalid CMD user role",
			func() fields {
				return fields{
					userRepo: &user.MockRepository{},
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Role(5)},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not update user, invalid role passed - 5", err.Error(), i)
				return true
			},
		},
		{
			"Can not get user from DB",
			func() fields {
				usrRepo := user.MockRepository{}
				usrRepo.On("Get", "testID").Return(nil, errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Admin, ID: "testID"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Can not check if new email is free",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				usrRepo.On("GetByEmail", "updated@test.com").Return(nil, errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Admin, ID: "testID", Email: "updated@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Email is not free",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				usrRepo.On("GetByEmail", "updated@test.com").Return(nil, nil)
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Admin, ID: "testID", Email: "updated@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "email updated@test.com already in use", err.Error(), i)
				return true
			},
		},
		{
			"Attempt to change user Role not from admin request",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Admin, ID: "testID", Email: "test@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not update user, attempt to change the role from not admin cmd", err.Error(), i)
				return true
			},
		},
		{
			"Error on passwd generation on not empty passwd",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				cipher := MockCipher{}
				cipher.On("GenerateHash", "passwd").Return("", errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &cipher,
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", Password: "passwd", Email: "test@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Error on applying changes",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				usrRepo.On("GetByEmail", "notValid").Return(nil, user.ErrNotFound)
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", Email: "notValid"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "email is not valid"), i)
				return true
			},
		},
		{
			"Error on changes saving",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				usrRepo.On("Update", mock.AnythingOfType("*user.User")).Return(errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", Email: "test@test.com", Name: "test1"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "testErr"), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateUserHandler{
				userRepo: tt.fieldsFn().userRepo,
				cipher:   tt.fieldsFn().cipher,
			}
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestUpdateUserHandler_Handle_PositiveCases(t *testing.T) {
	usrRepo := user.MockRepository{}
	usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
	assert.Nil(t, err)
	usrRepo.On("Get", "testID").Return(usr, nil)
	usrRepo.On("Update", mock.AnythingOfType("*user.User")).Return(nil)

	newPasswd := "passwd"
	newHash := "validPasswdHash"
	cipher := MockCipher{}
	cipher.On("GenerateHash", newPasswd).Return(newHash, nil)

	newEmail := "new@email.com"
	usrRepo.On("GetByEmail", newEmail).Return(nil, user.ErrNotFound)

	newRole := user.Admin
	cmd := UpdateUser{
		ID:         "testID",
		Name:       "newName",
		Email:      "new@email.com",
		Password:   newPasswd,
		Role:       newRole,
		IsAdminCMD: true,
	}

	handler := NewUpdateUserHandler(&usrRepo, &cipher)
	assert.Nil(t, handler.Handle(cmd))

	updatedUsr := usrRepo.Calls[2].Arguments[0].(*user.User)
	data := updatedUsr.ToMap()

	assert.Equal(t, cmd.Name, data["name"])
	assert.Equal(t, cmd.Email, data["email"])
	assert.Equal(t, newHash, data["password"])
	assert.Equal(t, int(cmd.Role), data["role"])
}

func TestUpdateUserHandler_emailFree(t *testing.T) {
	type fields struct {
		userRepo user.Repository
	}
	type args struct {
		cmd UpdateUser
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     bool
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on DB query",
			func() fields {
				repo := user.MockRepository{}
				repo.On("GetByEmail", "test@email.com").Return(nil, errors.New("testErr"))
				return fields{userRepo: &repo}
			},
			args{cmd: UpdateUser{Email: "test@email.com"}},
			false,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, err.Error(), "testErr", i)
				return true
			},
		},
		{
			"User not found - email is free",
			func() fields {
				repo := user.MockRepository{}
				repo.On("GetByEmail", "test@email.com").Return(nil, user.ErrNotFound)
				return fields{userRepo: &repo}
			},
			args{cmd: UpdateUser{Email: "test@email.com"}},
			true,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return true
			},
		},
	}
	cipher := MockCipher{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateUserHandler{
				userRepo: tt.fieldsFn().userRepo,
				cipher:   &cipher,
			}
			got, err := h.emailFree(tt.args.cmd)
			if !tt.wantErr(t, err, fmt.Sprintf("emailFree(%v)", tt.args.cmd)) {
				return
			}
			assert.Equalf(t, tt.want, got, "emailFree(%v)", tt.args.cmd)
		})
	}
}
