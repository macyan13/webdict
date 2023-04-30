package command

import "github.com/macyan13/webdict/backend/pkg/app/domain/translation"

type Opts struct {
	SupportedLanguages []translation.Lang
}
