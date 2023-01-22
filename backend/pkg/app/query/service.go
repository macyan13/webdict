package query

import (
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
