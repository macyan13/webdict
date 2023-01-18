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

func TestAddUserHandler_Handle_NegativeCases(t *testing.T) {
	type fields struct {
		userRepo user.Repository
		cipher   Cipher
	}
	type args struct {
		cmd AddUser
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error during passwd hash generation",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "testPwd").Return("", errors.New("testErr"))
				return fields{
					userRepo: &user.MockRepository{},
					cipher:   &cipher,
				}
			},
			args{cmd: AddUser{
				Name:     "testName",
				Email:    "test@email.com",
				Password: "testPwd",
				Role:     user.Admin,
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"User repo returns error on validation",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "testPwd").Return("hashedPwd", nil)
				userRepo := user.MockRepository{}
				userRepo.On("Exist", "test@email.com").Return(false, errors.New("testErr"))
				return fields{
					userRepo: &userRepo,
					cipher:   &cipher,
				}
			},
			args{cmd: AddUser{
				Name:     "testName",
				Email:    "test@email.com",
				Password: "testPwd",
				Role:     user.Admin,
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"User with such email already exist",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "testPwd").Return("hashedPwd", nil)
				userRepo := user.MockRepository{}
				userRepo.On("Exist", "test@email.com").Return(true, nil)
				return fields{
					userRepo: &userRepo,
					cipher:   &cipher,
				}
			},
			args{cmd: AddUser{
				Name:     "testName",
				Email:    "test@email.com",
				Password: "testPwd",
				Role:     user.Admin,
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not create new user, a user with passed email: test@email.com already exists", err.Error(), i)
				return true
			},
		},
		{
			"Error on user creation",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "testPwd").Return("hashedPasswd", nil)
				userRepo := user.MockRepository{}
				userRepo.On("Exist", "test@email.com").Return(false, nil)
				return fields{
					userRepo: &userRepo,
					cipher:   &cipher,
				}
			},
			args{cmd: AddUser{
				Name:     "n",
				Email:    "test@email.com",
				Password: "testPwd",
				Role:     user.Admin,
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "name must contain at least 2 characters, 1 passed (n)"), i)
				return true
			},
		},
		{
			"User repo returns error during user creation",
			func() fields {
				cipher := MockCipher{}
				cipher.On("GenerateHash", "testPwd").Return("hashedPwd", nil)
				userRepo := user.MockRepository{}
				userRepo.On("Exist", "test@email.com").Return(false, nil)
				userRepo.On("Create", mock.AnythingOfType("*user.User")).Return(errors.New("testErr"))
				return fields{
					userRepo: &userRepo,
					cipher:   &cipher,
				}
			},
			args{cmd: AddUser{
				Name:     "testName",
				Email:    "test@email.com",
				Password: "testPwd",
				Role:     user.Admin,
			}},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFn()
			h := NewAddUserHandler(fields.userRepo, fields.cipher)
			id, err := h.Handle(tt.args.cmd)
			assert.Equal(t, "", id)
			tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.cmd))
		})
	}
}

func TestAddUserHandler_Handle_PositiveCase(t *testing.T) {
	name := "testName"
	email := "test@email.com"
	pwd := "testPwd"
	hashedPwd := "hashedPwd"

	cipher := MockCipher{}
	cipher.On("GenerateHash", pwd).Return(hashedPwd, nil)
	userRepo := user.MockRepository{}
	userRepo.On("Exist", email).Return(false, nil)
	userRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)

	cmd := AddUser{
		Name:     name,
		Email:    email,
		Password: pwd,
		Role:     user.Admin,
	}

	handler := NewAddUserHandler(&userRepo, &cipher)

	id, err := handler.Handle(cmd)
	assert.Nil(t, err)

	createdUser := userRepo.Calls[1].Arguments[0].(*user.User)
	data := createdUser.ToMap()

	assert.Equal(t, createdUser.ID(), id)
	assert.Equal(t, cmd.Email, data["email"])
	assert.Equal(t, cmd.Name, data["name"])
	assert.Equal(t, hashedPwd, data["password"])
	assert.Equal(t, int(cmd.Role), data["role"])
}
