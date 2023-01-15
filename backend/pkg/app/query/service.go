package query

import (
	"github.com/microcosm-cc/bluemonday"
	"text/template"
)

var viewSanitizer *sanitizer

func init() {
	viewSanitizer = &sanitizer{policy: bluemonday.UGCPolicy()}
}

type sanitizer struct {
	policy *bluemonday.Policy
}

func (s sanitizer) Sanitize(input string) string {
	return s.policy.Sanitize(input)
}

func (s sanitizer) SanitizeAndEscape(input string) string {
	clean := s.policy.Sanitize(input)
	return template.HTMLEscapeString(clean)
}
