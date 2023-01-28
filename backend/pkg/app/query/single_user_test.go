package query

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingleUserHandler_Handle(t *testing.T) {
	type fields struct {
		userRepo   UserViewRepository
		strictSntz *strictSanitizer
	}
	type args struct {
		cmd SingleUser
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     UserView
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			"Error on getting users from DB",
			func() fields {
				userRepo := MockUserViewRepository{}
				userRepo.On("GetView", "testID").Return(UserView{}, errors.New("testErr"))
				return fields{userRepo: &userRepo}
			},
			args{
				cmd: SingleUser{
					ID: "testID",
				},
			},
			UserView{},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
		},
		{
			"Sanitize is called",
			func() fields {
				userRepo := MockUserViewRepository{}
				usr := UserView{
					ID:    "testId",
					Name:  `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">TestName<a>`,
					Email: `<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">TestEmail<a>`,
					Role:  2,
				}
				userRepo.On("GetView", "testID").Return(usr, nil)
				return fields{userRepo: &userRepo}
			},
			args{
				cmd: SingleUser{
					ID: "testID",
				},
			},
			UserView{
				ID:    "testId",
				Name:  "TestName",
				Email: "TestEmail",
				Role:  2,
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	s := newStrictSanitizer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := SingleUserHandler{
				userRepo:   tt.fieldsFn().userRepo,
				strictSntz: s,
			}
			got, err := h.Handle(tt.args.cmd)
			if !tt.wantErr(t, err, fmt.Sprintf("Handle(%v)", tt.args.cmd)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Handle(%v)", tt.args.cmd)
		})
	}
}
