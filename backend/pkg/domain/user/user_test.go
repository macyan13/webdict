package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
		role     Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Name is too short",
			args{
				"t",
				"test@test.com",
				"12345678",
				Admin,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not create new user, the name must contain at least 3 character", err.Error(), i)
				return true
			},
		},
		{
			"Password is too short",
			args{
				"tes",
				"test@test.com",
				"1234567",
				Admin,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not create new user, the password must contain at least 3 character", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			args{
				"test",
				"test@test.com",
				"12345678",
				Admin,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email, tt.args.password, tt.args.role)
			if tt.wantErr(t, err, fmt.Sprintf("NewUser(%v, %v, %v)", tt.args.name, tt.args.email, tt.args.password)) {
				return
			}
			assert.Equal(t, got.name, tt.args.name)
			assert.Equal(t, got.password, tt.args.password)
			assert.Equal(t, got.email, tt.args.email)
			assert.Equal(t, got.role, tt.args.role)
		})
	}
}

func TestUnmarshalFromDB(t *testing.T) {
	user := User{
		id:       "testId",
		name:     "testName",
		email:    "testEmail",
		password: "testPassword",
		role:     0,
	}

	assert.Equal(t, &user, UnmarshalFromDB(user.id, user.name, user.email, user.password, int(user.role)))
}
