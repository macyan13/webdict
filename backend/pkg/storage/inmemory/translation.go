package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"time"
)

type TranslationRepo struct {
	tagRepo TagRepo
	storage map[string]translation.Translation
}

func NewTranslationRepository(tagRepo TagRepo) *TranslationRepo {
	return &TranslationRepo{
		storage: map[string]translation.Translation{},
		tagRepo: tagRepo,
	}
}

func (r *TranslationRepo) Update(translation translation.Translation) error {
	r.storage[translation.Id()] = translation
	return nil
}

func (r *TranslationRepo) Get(id, authorId string) (translation.Translation, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorId() == authorId {
		return t, nil
	}

	return translation.Translation{}, translation.NotFoundErr
}

func (r *TranslationRepo) Delete(id, authorId string) error {
	t, ok := r.storage[id]

	if ok && t.AuthorId() == authorId {
		delete(r.storage, id)
		return nil
	}

	return fmt.Errorf("not found")
}

func (r *TranslationRepo) Create(translation translation.Translation) error {
	r.storage[translation.Id()] = translation
	return nil
}

func (r *TranslationRepo) GetView(id, authorId string) (query.TranslationView, error) {
	for _, t := range r.storage {

		if t.AuthorId() == authorId && t.Id() == id {
			translationData := t.ToMap()
			tagViews, err := r.tagRepo.GetViews(translationData["tagIds"].([]string), authorId)

			if err != nil {
				return query.TranslationView{}, err
			}
			return query.TranslationView{
				Id:            t.Id(),
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

func (r *TranslationRepo) GetLastViews(authorId string, limit int) ([]query.TranslationView, error) {
	results := make([]query.TranslationView, 0)
	counter := 0

	for _, t := range r.storage {
		if t.AuthorId() == authorId && counter < limit {
			translationData := t.ToMap()
			tagViews, err := r.tagRepo.GetViews(translationData["tagIds"].([]string), authorId)

			if err != nil {
				return nil, err
			}

			results = append(results, query.TranslationView{
				Id:            t.Id(),
				CreatedAd:     translationData["createdAt"].(time.Time),
				Transcription: translationData["transcription"].(string),
				Translation:   translationData["translation"].(string),
				Text:          translationData["text"].(string),
				Example:       translationData["example"].(string),
				Tags:          tagViews,
			})

			counter++
		}
	}

	return results, nil
}
