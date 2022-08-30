package repository

import (
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
)

type translationRepo struct {
	storage map[string]translation.Translation
}

func NewTranslationRepository() *translationRepo {
	return &translationRepo{
		storage: map[string]translation.Translation{},
	}
}

func (r *translationRepo) Save(translation translation.Translation) error {
	r.storage[translation.Id] = translation
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

func (r *translationRepo) Delete(id string) error {
	delete(r.storage, id)
	return nil
}
