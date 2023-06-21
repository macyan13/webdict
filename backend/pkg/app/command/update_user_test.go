package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
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
			"Error on new passwd generation",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				cipher := MockCipher{}
				cipher.On("GenerateHash", "newPasswd").Return("", errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &cipher,
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", Password: "newPasswd", Email: "test@test.com"},
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
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", Email: "notValid"},
			},
			assert.Error,
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
	currentPasswdHash := "testPasswd"
	usrRepo := user.MockRepository{}
	usr, err := user.NewUser("test", "test@test.com", currentPasswdHash, user.Author)
	assert.Nil(t, err)
	usrRepo.On("Get", "testID").Return(usr, nil)
	usrRepo.On("Update", mock.AnythingOfType("*user.User")).Return(nil)

	newPasswd := "passwd"
	newHash := "validPasswdHash"
	ID := "testID"
	cipher := MockCipher{}
	cipher.On("GenerateHash", newPasswd).Return(newHash, nil)

	newRole := user.Admin
	cmd := UpdateUser{
		ID:       ID,
		Name:     "newName",
		Email:    "new@email.com",
		Password: newPasswd,
		Role:     newRole,
	}

	handler := NewUpdateUserHandler(&usrRepo, &cipher)
	assert.Nil(t, handler.Handle(cmd))

	updatedUsr := usrRepo.Calls[1].Arguments[0].(*user.User)
	data := updatedUsr.ToMap()

	assert.Equal(t, cmd.Name, data["name"])
	assert.Equal(t, cmd.Email, data["email"])
	assert.Equal(t, newHash, data["password"])
	assert.Equal(t, int(cmd.Role), data["role"])
	assert.Equal(t, usr.DefaultLangID(), data["defaultLangID"])
}

func TestUpdateUserHandler_processPasswd(t *testing.T) {
	type fields struct {
		cipher Cipher
	}
	type args struct {
		cmd      UpdateUser
		userHash string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Password is not changed",
			func() fields {
				return fields{}
			},
			args{
				cmd:      UpdateUser{},
				userHash: "oldHash",
			},
			"oldHash",
			assert.NoError,
		},
		{
			"Password is changed",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "newPassword").Return("newHash", nil)
				return fields{cipher: &cipher}
			},
			args{
				cmd:      UpdateUser{Password: "newPassword"},
				userHash: "oldHash",
			},
			"newHash",
			assert.NoError,
		},
		{
			"Error on generating hash",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "newPassword").Return("", errors.New("test"))
				return fields{cipher: &cipher}
			},
			args{
				cmd:      UpdateUser{Password: "newPassword"},
				userHash: "oldHash",
			},
			"",
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateUserHandler{
				cipher: tt.fieldsFn().cipher,
			}
			got, err := h.processPasswd(tt.args.cmd, tt.args.userHash)
			if !tt.wantErr(t, err, fmt.Sprintf("processPasswd(%v, %v)", tt.args.cmd, tt.args.userHash)) {
				return
			}
			assert.Equalf(t, tt.want, got, "processPasswd(%v, %v)", tt.args.cmd, tt.args.userHash)
		})
	}
}
