package query

import (
	"time"
)

type TranslationViewRepository interface {
	GetView(id, authorID string) (TranslationView, error)
	GetLastViewsByTags(authorID, langID string, pageSize, page int, tagIDs []string) (LastTranslationViews, error)
	GetRandomViews(authorID, langID string, tagIDs []string, limit int) (RandomViews, error)
	GetLastViewsBySourcePart(authorID, langID, sourcePart string, pageSize, page int) (LastTranslationViews, error)
	GetLastViewsByTargetPart(authorID, langID, targetPart string, pageSize, page int) (LastTranslationViews, error)
}

type LastTranslationViews struct {
	Views        []TranslationView
	TotalRecords int
}

type RandomViews struct {
	Views []TranslationView
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

type RoleView struct {
	ID      int
	Name    string
	IsAdmin bool
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
	ID   string
	Name string
}

func (v *TagView) sanitize(sanitizer *strictSanitizer) {
	v.Name = sanitizer.Sanitize(v.Name)
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
	Role        RoleView
	DefaultLang LangView
	ListOptions UserListOptionsView
}

type UserListOptionsView struct {
	HideTranscription bool
}

func (v *UserView) sanitize(sanitizer *strictSanitizer) {
	v.Name = sanitizer.Sanitize(v.Name)
	v.Email = sanitizer.Sanitize(v.Email)
	v.DefaultLang.sanitize(sanitizer)
}
