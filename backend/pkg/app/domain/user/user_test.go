package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
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
				assert.True(t, strings.Contains(err.Error(), "name must contain at least 2 characters, 1 passed (t)"), i)
				return true
			},
		},
		{
			"Name is too long",
			args{
				string(make([]rune, 31)),
				"test@test.com",
				"12345678",
				Admin,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "name max size is 30 characters, 31 passed"), i)
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
				assert.True(t, strings.Contains(err.Error(), "password must contain at least 8 character"), i)
				return true
			},
		},
		{
			"Invalid Email",
			args{
				"tes",
				"test.test.com",
				"1234567",
				Admin,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "email is not valid"), i)
				return true
			},
		},
		{
			"Invalid Role",
			args{
				"tes",
				"test.test.com",
				"1234567",
				Role(0),
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "invalid user role passed - 0"), i)
				return true
			},
		},
		{
			"Multiple errors",
			args{
				"tes",
				"test.test.com",
				"12367",
				Admin,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "email is not valid"), i)
				assert.True(t, strings.Contains(err.Error(), "password must contain at least 8 character"), i)
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
		id:            "testId",
		name:          "testName",
		email:         "testEmail",
		password:      "testPassword",
		role:          0,
		defaultLangID: "testLang",
	}

	assert.Equal(t, &user, UnmarshalFromDB(user.id, user.name, user.email, user.password, int(user.role), user.defaultLangID))
}

func TestRole_valid(t *testing.T) {
	tests := []struct {
		name string
		r    Role
		want bool
	}{
		{
			"Invalid role, value less than the actual min",
			Role(0),
			false,
		},
		{
			"Invalid role, value bigger than the actual min",
			Role(3),
			false,
		},
		{
			"Positive case",
			Admin,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.r.valid(), "valid()")
		})
	}
}

func TestUser_ApplyChanges(t *testing.T) {
	type fields struct {
		name     string
		email    string
		password string
		role     Role
	}
	type args struct {
		name          string
		email         string
		passwd        string
		role          Role
		defaultLangID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantFn func(t assert.TestingT, err error, usr *User, details string)
	}{
		{
			"Error on validation, changes should not be applied",
			fields{
				name:     "testName",
				email:    "test@mail.com",
				password: "testPasswd",
				role:     Admin,
			},
			args{
				name:   "name",
				email:  "invalidEmail",
				passwd: "testPasswd",
				role:   Author,
			},
			func(t assert.TestingT, err error, usr *User, details string) {
				assert.True(t, strings.Contains(err.Error(), "email is not valid"), details)
				assert.Equal(t, "testName", usr.name)
				assert.Equal(t, "test@mail.com", usr.email)
				assert.Equal(t, "testPasswd", usr.password)
				assert.Equal(t, Admin, usr.role)
			},
		},
		{
			"Applied changes",
			fields{
				name:     "testName",
				email:    "test@mail.com",
				password: "testPasswd",
				role:     Admin,
			},
			args{
				name:          "name",
				email:         "updated@email.com",
				passwd:        "updatedPasswd",
				role:          Author,
				defaultLangID: "langID",
			},
			func(t assert.TestingT, err error, usr *User, details string) {
				assert.Nil(t, err, details)
				assert.Equal(t, "name", usr.name)
				assert.Equal(t, "updated@email.com", usr.email)
				assert.Equal(t, "updatedPasswd", usr.password)
				assert.Equal(t, Author, usr.role)
				assert.Equal(t, "langID", usr.defaultLangID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				name:     tt.fields.name,
				email:    tt.fields.email,
				password: tt.fields.password,
				role:     tt.fields.role,
			}
			tt.wantFn(t, u.ApplyChanges(tt.args.name, tt.args.email, tt.args.passwd, tt.args.role, tt.args.defaultLangID), u, fmt.Sprintf("ApplyChanges(%v, %v, %v, %v, %v)", tt.args.name, tt.args.email, tt.args.passwd, tt.args.role, tt.args.defaultLangID))
		})
	}
}
