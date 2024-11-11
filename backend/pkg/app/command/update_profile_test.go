package command

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateProfileHandler_processPasswd(t *testing.T) {
	type fields struct {
		cipher Cipher
	}
	type args struct {
		cmd      UpdateProfile
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
			"New password is not set",
			func() fields {
				return fields{cipher: &MockCipher{}}
			},
			args{
				cmd:      UpdateProfile{},
				userHash: "test",
			},
			"test",
			assert.NoError,
		},
		{
			"User password and current one are different",
			func() fields {
				cipher := MockCipher{}
				cipher.On("ComparePasswords", "userHash", "current").Return(false)
				return fields{cipher: &cipher}
			},
			args{
				cmd:      UpdateProfile{NewPassword: "test", CurrentPassword: "current"},
				userHash: "userHash",
			},
			"",
			assert.Error,
		},
		{
			"Error on password generation",
			func() fields {
				cipher := MockCipher{}
				cipher.On("ComparePasswords", "userHash", "current").Return(true)
				cipher.On("GenerateHash", "test").Return("", errors.New("test"))
				return fields{cipher: &cipher}
			},
			args{
				cmd:      UpdateProfile{NewPassword: "test", CurrentPassword: "current"},
				userHash: "userHash",
			},
			"",
			assert.Error,
		},
		{
			"New password is generated",
			func() fields {
				cipher := MockCipher{}
				cipher.On("ComparePasswords", "userHash", "current").Return(true)
				cipher.On("GenerateHash", "test").Return("newHash", nil)
				return fields{cipher: &cipher}
			},
			args{
				cmd:      UpdateProfile{NewPassword: "test", CurrentPassword: "current"},
				userHash: "userHash",
			},
			"newHash",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateProfileHandler{
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

func TestUpdateProfileHandler_validate(t *testing.T) {
	type fields struct {
		langRepo lang.Repository
	}
	type args struct {
		cmd UpdateProfile
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Default lang is not set",
			func() fields {
				return fields{langRepo: &lang.MockRepository{}}
			},
			args{cmd: UpdateProfile{}},
			assert.NoError,
		},
		{
			"Error on getting lang from DB by id",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "userID").Return(false, errors.New("test"))
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateProfile{ID: "userID", DefaultLangID: "langID"}},
			assert.Error,
		},
		{
			"Lang does not exist",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "userID").Return(false, nil)
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateProfile{ID: "userID", DefaultLangID: "langID"}},
			assert.Error,
		},
		{
			"Lang exists",
			func() fields {
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "langID", "userID").Return(true, nil)
				return fields{langRepo: &langRepo}
			},
			args{cmd: UpdateProfile{ID: "userID", DefaultLangID: "langID"}},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateProfileHandler{
				langRepo: tt.fieldsFn().langRepo,
			}
			tt.wantErr(t, h.validate(tt.args.cmd), fmt.Sprintf("validate(%v)", tt.args.cmd))
		})
	}
}

func TestUpdateProfileHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		cipher   Cipher
		langRepo lang.Repository
	}
	type args struct {
		cmd UpdateProfile
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
				cmd: UpdateProfile{ID: "testID"},
			},
			assert.Error,
		},
		{
			"Validate - Error on getting Lang from DB",
			func() fields {
				usrRepo := user.MockRepository{}
				usr, err := user.NewUser("test", "test@test.com", "testPasswd", user.Author)
				assert.Nil(t, err)
				usrRepo.On("Get", "testID").Return(usr, nil)
				langRepo := lang.MockRepository{}
				langRepo.On("Exist", "testLangID", "testID").Return(false, errors.New("testErr"))
				return fields{
					userRepo: &usrRepo,
					cipher:   &MockCipher{},
					langRepo: &langRepo,
				}
			},
			args{
				cmd: UpdateProfile{ID: "testID", CurrentPassword: "passwd", NewPassword: "newPasswd", Email: "test@test.com", DefaultLangID: "testLangID"},
			},
			assert.Error,
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
				cmd: UpdateProfile{ID: "testID", CurrentPassword: "passwd", NewPassword: "test", Email: "test@test.com"},
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "current password is not valid", err.Error(), i)
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
				cmd: UpdateProfile{ID: "testID", Email: "notValid"},
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
				cmd: UpdateProfile{ID: "testID", Email: "test@test.com", Name: "test1"},
			},
			assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := UpdateProfileHandler{
				userRepo: tt.fieldsFn().userRepo,
				cipher:   tt.fieldsFn().cipher,
				langRepo: tt.fieldsFn().langRepo,
			}
			tt.wantErr(t, h.Handle(tt.args.cmd), fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestUpdateProfileHandler_Handle_PositiveCases(t *testing.T) {
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
	langID := "testLang"
	langRepo.On("Exist", langID, ID).Return(true, nil)

	cmd := UpdateProfile{
		ID:              ID,
		Name:            "newName",
		Email:           "new@email.com",
		CurrentPassword: currentPasswd,
		NewPassword:     newPasswd,
		DefaultLangID:   langID,
		ListOptions:     user.NewListOptions(true),
	}

	handler := NewUpdateProfileHandler(&usrRepo, &cipher, &langRepo)
	assert.Nil(t, handler.Handle(cmd))

	updatedUsr := usrRepo.Calls[1].Arguments[0].(*user.User)
	data := updatedUsr.ToMap()
	listData := updatedUsr.ListOptions()

	assert.Equal(t, cmd.Name, data["name"])
	assert.Equal(t, cmd.Email, data["email"])
	assert.Equal(t, newHash, data["password"])
	assert.Equal(t, langID, data["defaultLangID"])
	assert.Equal(t, true, listData.ToMap()["hideTranscription"])
}
