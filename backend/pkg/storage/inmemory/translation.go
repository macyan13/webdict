package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"time"
)

type TranslationRepo struct {
	tagRepo TagRepo
	storage map[string]*translation.Translation
}

func NewTranslationRepository(tagRepo TagRepo) *TranslationRepo {
	return &TranslationRepo{
		storage: map[string]*translation.Translation{},
		tagRepo: tagRepo,
	}
}

func (r *TranslationRepo) Update(t translation.Translation) error {
	r.storage[t.ID()] = &t
	return nil
}

func (r *TranslationRepo) Get(id, authorID string) (translation.Translation, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		return *t, nil
	}

	return translation.Translation{}, translation.ErrNotFound
}

func (r *TranslationRepo) Delete(id, authorID string) error {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		delete(r.storage, id)
		return nil
	}

	return fmt.Errorf("not found")
}

func (r *TranslationRepo) Create(t translation.Translation) error {
	r.storage[t.ID()] = &t
	return nil
}

func (r *TranslationRepo) GetView(id, authorID string) (query.TranslationView, error) {
	for _, t := range r.storage {

		if t.AuthorID() == authorID && t.ID() == id {
			translationData := t.ToMap()
			tagViews, err := r.tagRepo.GetViews(translationData["tagIDs"].([]string), authorID)

			if err != nil {
				return query.TranslationView{}, err
			}
			return query.TranslationView{
				ID:            t.ID(),
				CreatedAd:     translationData["createdAt"].(time.Time),
				Transcription: translationData["transcription"].(string),
				Translation:   translationData["translation"].(string),
				Text:          translationData["text"].(string),
				Example:       translationData["example"].(string),
				Tags:          tagViews,
			}, nil
		}
	}

	return query.TranslationView{}, fmt.Errorf("not found")
}

func (r *TranslationRepo) GetLastViews(authorID string, limit int) ([]query.TranslationView, error) {
	results := make([]query.TranslationView, 0)
	counter := 0

	for _, t := range r.storage {
		if t.AuthorID() != authorID || counter >= limit {
			continue
		}

		translationData := t.ToMap()
		tagViews, err := r.tagRepo.GetViews(translationData["tagIDs"].([]string), authorID)

		if err != nil {
			return nil, err
		}

		results = append(results, query.TranslationView{
			ID:            t.ID(),
			CreatedAd:     translationData["createdAt"].(time.Time),
			Transcription: translationData["transcription"].(string),
			Translation:   translationData["translation"].(string),
			Text:          translationData["text"].(string),
			Example:       translationData["example"].(string),
			Tags:          tagViews,
		})

		counter++
	}

	return results, nil
}
