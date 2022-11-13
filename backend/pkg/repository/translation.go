package repository

import (
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type translationRepo struct {
	tagRepo tagRepo
	storage map[string]translation.Translation
}

func NewTranslationRepository(tagRepo tagRepo) *translationRepo {
	return &translationRepo{
		storage: map[string]translation.Translation{},
		tagRepo: tagRepo,
	}
}

func (r *translationRepo) Save(translation translation.Translation) error {
	r.storage[translation.Id()] = translation
	return nil
}

func (r *translationRepo) Get() []translation.Translation {
	result := make([]translation.Translation, 0, len(r.storage))

	for _, translation := range r.storage {
		result = append(result, translation)
	}

	return result
}

func (r *translationRepo) GetById(id string) *translation.Translation {
	translation, ok := r.storage[id]

	if ok {
		return &translation
	}

	return nil
}

func (r *translationRepo) Delete(translation translation.Translation) error {
	delete(r.storage, translation.Id())
	return nil
}

func (r *translationRepo) GetLastTranslations(authorId string, limit int) []query.Translation {
	results := make([]query.Translation, 0)
	counter := 0

	for _, t := range r.storage {
		if t.AuthorId() == authorId && counter < limit {
			results = append(results, query.Translation{
				Id:            t.Id(),
				CreatedAd:     t.CreatedAt(),
				Transcription: t.Transcription(),
				Translation:   t.Translation(),
				Text:          t.Text(),
				Example:       t.Example(),
				Tags:          r.tagRepo.getTagsById(t.TagIds()),
			})

			counter++
		}
	}

	return results
}

func (r *translationRepo) GetTranslation(id, authorId string) *query.Translation {
	for _, t := range r.storage {
		if t.AuthorId() == authorId && t.Id() == id {
			return &query.Translation{
				Id:            t.Id(),
				CreatedAd:     t.CreatedAt(),
				Transcription: t.Transcription(),
				Translation:   t.Translation(),
				Text:          t.Text(),
				Example:       t.Example(),
				Tags:          r.tagRepo.getTagsById(t.TagIds()),
			}
		}
	}

	return nil
}
