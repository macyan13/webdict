package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
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
		langRepo lang.Repository
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
			"Validate - Error on getting Lang from DB",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				langRepo := lang.MockRepository{}
				langRepo.On("Get", "testLangID", "testID").Return(nil, errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
					langRepo: &langRepo,
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", CurrentPassword: "passwd", NewPassword: "newPasswd", Email: "test@test.com", DefaultLangID: "testLangID"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Current passwd hash and user hash are not equal",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				cipher := MockCipher{}
				cipher.On("ComparePasswords", "testPasswd", "passwd").Return(false)
				return fields{
					userRepo: &usrRepo,
					cipher:   &cipher,
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", CurrentPassword: "passwd", NewPassword: "test", Email: "test@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "current password is not valid", err.Error(), i)
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
				cipher.On("ComparePasswords", "testPasswd", "passwd").Return(true)
				cipher.On("GenerateHash", "newPasswd").Return("", errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &cipher,
				}
			},
			args{
				cmd: UpdateUser{Role: user.Author, ID: "testID", CurrentPassword: "passwd", NewPassword: "newPasswd", Email: "test@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
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
				assert.Equal(t, "attempt to change the role from not admin cmd", err.Error(), i)
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
				langRepo: tt.fieldsFn().langRepo,
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
	currentPasswd := "currentPasswd"
	ID := "testID"
	cipher := MockCipher{}
	cipher.On("ComparePasswords", currentPasswdHash, currentPasswd).Return(true)
	cipher.On("GenerateHash", newPasswd).Return(newHash, nil)

	langRepo := lang.MockRepository{}
	ln, err := lang.NewLang("EN", ID)
	assert.Nil(t, err)
	langRepo.On("Get", ln.ID(), ID).Return(ln, nil)

	newRole := user.Admin
	cmd := UpdateUser{
		ID:              ID,
		Name:            "newName",
		Email:           "new@email.com",
		CurrentPassword: currentPasswd,
		NewPassword:     newPasswd,
		Role:            newRole,
		IsAdminCMD:      true,
		DefaultLangID:   ln.ID(),
	}

	handler := NewUpdateUserHandler(&usrRepo, &cipher, &langRepo)
	assert.Nil(t, handler.Handle(cmd))

	updatedUsr := usrRepo.Calls[1].Arguments[0].(*user.User)
	data := updatedUsr.ToMap()

	assert.Equal(t, cmd.Name, data["name"])
	assert.Equal(t, cmd.Email, data["email"])
	assert.Equal(t, newHash, data["password"])
	assert.Equal(t, int(cmd.Role), data["role"])
	assert.Equal(t, ln.ID(), data["defaultLangID"])
}

func TestUpdateUserHandler_processRole(t *testing.T) {
	type args struct {
		cmd UpdateUser
		usr *user.User
	}
	tests := []struct {
		name    string
		argsFn  func() args
		want    user.Role
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Admin command",
			func() args {
				return args{
					cmd: UpdateUser{Role: user.Admin, IsAdminCMD: true},
					usr: nil,
				}
			},
			user.Admin,
			assert.NoError,
		},
		{
			"Not admin command, roles are different",
			func() args {
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				return args{
					cmd: UpdateUser{Role: user.Admin},
					usr: usr,
				}
			},
			0,
			assert.Error,
		},
		{
			"Not admin command, role is not set, return user role",
			func() args {
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				return args{
					cmd: UpdateUser{Role: user.NotSet},
					usr: usr,
				}
			},
			user.Author,
			assert.NoError,
		},
		{
			"Not admin command, role is set, return user role",
			func() args {
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				return args{
					cmd: UpdateUser{Role: user.Author},
					usr: usr,
				}
			},
			user.Author,
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateUserHandler{}
			got, err := h.processRole(tt.argsFn().cmd, tt.argsFn().usr)
			if !tt.wantErr(t, err, fmt.Sprintf("processRole(%v, %v)", tt.argsFn().cmd, tt.argsFn().usr)) {
				return
			}
			assert.Equalf(t, tt.want, got, "processRole(%v, %v)", tt.argsFn().cmd, tt.argsFn().usr)
		})
	}
}
