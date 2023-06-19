package query

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoleConverter_RoleToView(t *testing.T) {
	type args struct {
		role user.Role
	}
	tests := []struct {
		name    string
		args    args
		want    RoleView
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Error on not mapped role",
			args{role: user.Role(3)},
			RoleView{},
			assert.Error,
		},
		{
			"Author",
			args{role: user.Role(2)},
			RoleView{Name: "User", ID: 2, IsAdmin: false},
			assert.NoError,
		},
		{
			"Admin",
			args{role: user.Role(1)},
			RoleView{Name: "Admin", ID: 1, IsAdmin: true},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewRoleMapper()
			got, err := m.RoleToView(tt.args.role)
			if !tt.wantErr(t, err, fmt.Sprintf("RoleToView(%v)", tt.args.role)) {
				return
			}
			assert.Equalf(t, tt.want, got, "RoleToView(%v)", tt.args.role)
		})
	}
}

func Test_richTextSanitizer_SanitizeAndEscape(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Malicious content",
			args{input: `<a onblur="alert(secret)" href="http://www.test.com">test</a>`},
			`&lt;a href=&#34;http://www.test.com&#34; rel=&#34;nofollow&#34;&gt;test&lt;/a&gt;`,
		},
		{
			"Safe content",
			args{input: "just string"},
			"just string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newRichTextSanitizer()
			assert.Equalf(t, tt.want, s.SanitizeAndEscape(tt.args.input), "SanitizeAndEscape(%v)", tt.args.input)
		})
	}
}
