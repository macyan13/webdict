package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"sort"
	"time"
)

type TranslationRepo struct {
	tagRepo  TagRepo
	langRepo LangRepo
	storage  map[string]*translation.Translation
}

func NewTranslationRepository(tagRepo TagRepo, langRepo LangRepo) *TranslationRepo {
	return &TranslationRepo{
		tagRepo:  tagRepo,
		langRepo: langRepo,
		storage:  map[string]*translation.Translation{},
	}
}

func (r *TranslationRepo) Update(t *translation.Translation) error {
	r.storage[t.ID()] = t
	return nil
}

func (r *TranslationRepo) Get(id, authorID string) (*translation.Translation, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		return t, nil
	}

	return nil, translation.ErrNotFound
}

func (r *TranslationRepo) Delete(id, authorID string) error {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		delete(r.storage, id)
		return nil
	}

	return fmt.Errorf("not found")
}

func (r *TranslationRepo) Create(t *translation.Translation) error {
	r.storage[t.ID()] = t
	return nil
}

func (r *TranslationRepo) ExistByLang(langID, authorID string) (bool, error) {
	for _, t := range r.storage {
		if t.AuthorID() != authorID || t.LangID() != langID {
			continue
		}

		return true, nil
	}

	return false, nil
}

func (r *TranslationRepo) ExistByTag(tagID, authorID string) (bool, error) {
	for _, t := range r.storage {
		if t.AuthorID() != authorID {
			continue
		}

		for _, tag := range t.ToMap()["tagIDs"].([]string) {
			if tag == tagID {
				return true, nil
			}
		}
	}

	return false, nil
}

func (r *TranslationRepo) GetLastViews(authorID, langID string, pageSize, page int, tagIds []string) (query.LastViews, error) {
	type mapItem struct {
		t         *translation.Translation
		createdAt time.Time
	}

	items := make([]mapItem, 0, len(r.storage))

	for _, v := range r.storage {
		if v.AuthorID() != authorID || v.LangID() != langID {
			continue
		}

		data := v.ToMap()

		if !r.containsAll(data["tagIDs"].([]string), tagIds) {
			continue
		}

		items = append(items, mapItem{
			t:         v,
			createdAt: data["createdAt"].(time.Time),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].createdAt.After(items[j].createdAt)
	})

	var offset int
	if page > 1 {
		offset = pageSize * page
	}

	if len(items) < offset && page != 1 {
		return query.LastViews{}, fmt.Errorf("can not get translations from DB")
	}

	views := make([]query.TranslationView, 0, len(items))

	i := 0
	limit := offset + pageSize
	for _, v := range items {
		if i >= offset && i < limit {
			view, err := r.translationToView(v.t)

			if err != nil {
				return query.LastViews{}, err
			}
			views = append(views, view)
		} else if i >= offset+limit {
			break
		}
		i++
	}

	return query.LastViews{
		Views:        views,
		TotalRecords: len(items),
	}, nil
}

func (r *TranslationRepo) containsAll(tags, searchTags []string) bool {
	for _, searchTag := range searchTags {
		found := false
		for _, tag := range tags {
			if searchTag == tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (r *TranslationRepo) GetView(id, authorID string) (query.TranslationView, error) {
	for _, t := range r.storage {

		if t.AuthorID() == authorID && t.ID() == id {
			return r.translationToView(t)
		}
	}

	return query.TranslationView{}, fmt.Errorf("not found")
}

func (r *TranslationRepo) translationToView(t *translation.Translation) (query.TranslationView, error) {
	translationData := t.ToMap()
	tagViews, err := r.tagRepo.GetViews(translationData["tagIDs"].([]string), translationData["authorID"].(string))
	if err != nil {
		return query.TranslationView{}, err
	}

	langView, err := r.langRepo.GetView(translationData["langID"].(string), translationData["authorID"].(string))
	if err != nil {
		return query.TranslationView{}, err
	}

	return query.TranslationView{
		ID:            t.ID(),
		CreatedAd:     translationData["createdAt"].(time.Time),
		Transcription: translationData["transcription"].(string),
		Target:        translationData["target"].(string),
		Source:        translationData["source"].(string),
		Example:       translationData["example"].(string),
		Tags:          tagViews,
		Lang:          langView,
	}, nil
}
