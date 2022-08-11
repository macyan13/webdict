package repository

import (
	"github.com/Yan-Matskevich/webdict/backend/pkg/domain"
)

type translationRepo struct {
	storage map[string]domain.Translation
}

func NewTranslationRepository() *translationRepo {
	return &translationRepo{
		storage: map[string]domain.Translation{},
	}
}

func (r translationRepo) Save(translation domain.Translation) error {
	r.storage[translation.Id] = translation
	return nil
}

func (r translationRepo) Get() []domain.Translation {
	result := make([]domain.Translation, 0, len(r.storage))

	for _, translation := range r.storage {
		result = append(result, translation)
	}

	return result
}

func (r translationRepo) GetById(id string) *domain.Translation {
	translation, ok := r.storage[id]

	if ok {
		return &translation
	}

	return nil
}

func (r translationRepo) Delete(id string) error {
	delete(r.storage, id)
	return nil
}
