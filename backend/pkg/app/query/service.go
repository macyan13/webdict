package query

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/microcosm-cc/bluemonday"
	"text/template"
)

type sanitizer struct {
	policy *bluemonday.Policy
}

func (s sanitizer) Sanitize(input string) string {
	return s.policy.Sanitize(input)
}

type strictSanitizer struct {
	sanitizer
}

func newStrictSanitizer() *strictSanitizer {
	return &strictSanitizer{sanitizer{policy: bluemonday.StrictPolicy()}}
}

type richTextSanitizer struct {
	sanitizer
}

func newRichTextSanitizer() *richTextSanitizer {
	return &richTextSanitizer{sanitizer{policy: bluemonday.UGCPolicy()}}
}

func (s *richTextSanitizer) SanitizeAndEscape(input string) string {
	clean := s.policy.Sanitize(input)
	return template.HTMLEscapeString(clean)
}

type RoleConverter struct {
	nameMap map[user.Role]string
}

func NewRoleMapper() *RoleConverter {
	return &RoleConverter{nameMap: map[user.Role]string{
		user.Admin:  "Admin",
		user.Author: "User",
	}}
}

func (m RoleConverter) RoleToView(role user.Role) (RoleView, error) {
	name, ok := m.nameMap[role]

	if !ok {
		return RoleView{}, fmt.Errorf("name mapping for role %v is not set", role)
	}

	return RoleView{
		ID:      int(role),
		Name:    name,
		IsAdmin: role == user.Admin,
	}, nil
}
