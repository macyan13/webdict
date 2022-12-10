package query

import "time"

type TranslationViewRepository interface {
	GetView(id, authorId string) (TranslationView, error)
	GetLastViews(authorId string, limit int) ([]TranslationView, error)
}

type TagViewRepository interface {
	GetAllViews(authorId string) ([]TagView, error)
	GetView(id, authorId string) (TagView, error)
	GetViews(ids []string, authorId string) ([]TagView, error)
}

type TranslationView struct {
	Id            string
	CreatedAd     time.Time
	Transcription string
	Translation   string
	Text          string
	Example       string
	Tags          []TagView
}

type TagView struct {
	Id  string
	Tag string
}
