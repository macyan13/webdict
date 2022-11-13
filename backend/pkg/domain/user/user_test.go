package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
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
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
			if tt.wantErr(t, err, fmt.Sprintf("NewUser(%v, %v, %v)", tt.args.name, tt.args.email, tt.args.password)) {
				return
			}
			assert.Equal(t, time.Now().Format("2006-01-02 15:04:05"), got.createdAt.Format("2006-01-02 15:04:05"))
			assert.Equal(t, got.name, tt.args.name)
			assert.Equal(t, got.password, tt.args.password)
			assert.Equal(t, got.email, tt.args.email)
		})
	}
}
