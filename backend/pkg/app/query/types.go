package query

import "time"

type TranslationViewRepository interface {
	GetView(id, authorID string) (TranslationView, error)
	GetLastViews(authorID string, limit int) ([]TranslationView, error)
}

type TagViewRepository interface {
	GetAllViews(authorID string) ([]TagView, error)
	GetView(id, authorID string) (TagView, error)
	GetViews(ids []string, authorID string) ([]TagView, error)
}

type TranslationView struct {
	ID            string
	CreatedAd     time.Time
	Transcription string
	Translation   string
	Text          string
	Example       string
	Tags          []TagView
}

type TagView struct {
	ID  string
	Tag string
}
