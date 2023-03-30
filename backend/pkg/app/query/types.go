package query

import "time"

type TranslationViewRepository interface {
	GetView(id, authorID string) (TranslationView, error)
	GetLastViews(authorID string, pageSize, page int, tagIds []string) (LastViews, error)
}

type LastViews struct {
	Views      []TranslationView
	TotalPages int
}

type TagViewRepository interface {
	GetAllViews(authorID string) ([]TagView, error)
	GetView(id, authorID string) (TagView, error)
	GetViews(ids []string, authorID string) ([]TagView, error)
}

type UserViewRepository interface {
	GetAllViews() ([]UserView, error)
	GetView(id string) (UserView, error)
}

type TranslationView struct {
	ID            string
	Text          string
	Transcription string
	Meaning       string
	Example       string
	Tags          []TagView
	CreatedAd     time.Time
}

func (v *TranslationView) sanitize(strictSntz *strictSanitizer, reachSntz *richTextSanitizer) {
	v.Text = reachSntz.SanitizeAndEscape(v.Text)
	v.Transcription = reachSntz.SanitizeAndEscape(v.Transcription)
	v.Meaning = reachSntz.Sanitize(v.Meaning)
	v.Example = reachSntz.Sanitize(v.Example)

	for i := range v.Tags {
		v.Tags[i].sanitize(strictSntz)
	}
}

type TagView struct {
	ID  string
	Tag string
}

func (v *TagView) sanitize(sanitizer *strictSanitizer) {
	v.Tag = sanitizer.Sanitize(v.Tag)
}

type UserView struct {
	ID    string
	Name  string
	Email string
	Role  int
}

func (v *UserView) sanitize(sanitizer *strictSanitizer) {
	v.Name = sanitizer.Sanitize(v.Name)
	v.Email = sanitizer.Sanitize(v.Email)
}
