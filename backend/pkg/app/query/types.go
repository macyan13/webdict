package query

import "time"

type TranslationViewRepository interface {
	GetView(id, authorID string) (TranslationView, error)
	GetLastViews(authorID, langID string, pageSize, page int, tagIds []string) (LastViews, error)
}

type LastViews struct {
	Views        []TranslationView
	TotalRecords int
}

type TagViewRepository interface {
	GetAllViews(authorID string) ([]TagView, error)
	GetView(id, authorID string) (TagView, error)
	GetViews(ids []string, authorID string) ([]TagView, error)
}

type LangViewRepository interface {
	GetAllViews(authorID string) ([]LangView, error)
	GetView(id, authorID string) (LangView, error)
}

type UserViewRepository interface {
	GetAllViews() ([]UserView, error)
	GetView(id string) (UserView, error)
}

type TranslationView struct {
	ID            string
	Source        string
	Transcription string
	Target        string
	Example       string
	Tags          []TagView
	CreatedAd     time.Time
	Lang          LangView
}

func (v *TranslationView) sanitize(strictSntz *strictSanitizer, reachSntz *richTextSanitizer) {
	v.Source = reachSntz.SanitizeAndEscape(v.Source)
	v.Transcription = reachSntz.SanitizeAndEscape(v.Transcription)
	v.Target = reachSntz.Sanitize(v.Target)
	v.Example = reachSntz.Sanitize(v.Example)
	v.Lang.sanitize(strictSntz)

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

type LangView struct {
	ID   string
	Name string
}

func (v *LangView) sanitize(sanitizer *strictSanitizer) {
	v.Name = sanitizer.Sanitize(v.Name)
}

type UserView struct {
	ID          string
	Name        string
	Email       string
	Role        int
	DefaultLang LangView
}

func (v *UserView) sanitize(sanitizer *strictSanitizer) {
	v.Name = sanitizer.Sanitize(v.Name)
	v.Email = sanitizer.Sanitize(v.Email)
	v.DefaultLang.sanitize(sanitizer)
}
