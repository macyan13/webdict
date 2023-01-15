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
	Text          string
	Transcription string
	Translation   string
	Example       string
	Tags          []TagView
	CreatedAd     time.Time
}

func (v *TranslationView) sanitise() {
	v.Text = viewSanitizer.SanitizeAndEscape(v.Text)
	v.Transcription = viewSanitizer.SanitizeAndEscape(v.Transcription)
	v.Translation = viewSanitizer.Sanitize(v.Translation)
	v.Example = viewSanitizer.Sanitize(v.Example)

	for i := range v.Tags {
		v.Tags[i].sanitise()
	}
}

type TagView struct {
	ID  string
	Tag string
}

func (v *TagView) sanitise() {
	v.Tag = viewSanitizer.SanitizeAndEscape(v.Tag)
}
