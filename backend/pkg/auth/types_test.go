package auth

import (
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_IsAdmin(t *testing.T) {
	type fields struct {
		Role user.Role
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"Case 1: Is Admin",
			fields{Role: user.Admin},
			true,
		},
		{
			"Case 1: Is Author",
			fields{Role: user.Author},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				Role: tt.fields.Role,
			}
			assert.Equalf(t, tt.want, u.IsAdmin(), "IsAdmin()")
		})
	}
}
