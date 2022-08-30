package translation

import (
	"errors"
	"fmt"
)

type Repository interface {
	Save(translation Translation) error
	GetById(id string) *Translation
	Get() []Translation
	Delete(id string) error
}

type Service struct {
	repository Repository
}

func NewTranslationService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateTranslation(request Request) error {
	translation := NewTranslation(request)
	return s.repository.Save(*translation)
}

func (s *Service) UpdateTranslation(id string, request Request) error {
	translation := s.repository.GetById(id)

	if translation == nil {
		return errors.New(fmt.Sprintf("Can not find translation by ID: %s", id))
	}

	translation.ApplyChanges(request)
	return s.repository.Save(*translation)
}

func (s *Service) GetTranslations() []Translation {
	return s.repository.Get()
}

func (s *Service) GetById(id string) *Translation {
	return s.repository.GetById(id)
}

func (s *Service) DeleteById(id string) error {
	return s.repository.Delete(id)
}
