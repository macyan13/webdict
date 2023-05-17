package cache

import "time"

type Opts struct {
	TagCacheTTL         time.Duration
	TranslationCacheTTL time.Duration
	LangCacheTTL        time.Duration
}
