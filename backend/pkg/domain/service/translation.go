package service

import (
	"errors"
	"fmt"
	"github.com/Yan-Matskevich/webdict/backend/pkg/domain"
)

type TranslationRepository interface {
	Save(translation domain.Translation) error
	GetById(id string) *domain.Translation
	Get() []domain.Translation
	Delete(id string) error
}

type TranslationService struct {
	repository TranslationRepository
}

func NewTranslationService(repository TranslationRepository) *TranslationService {
	return &TranslationService{
		repository: repository,
	}
}

func (s *TranslationService) CreateTranslation(request domain.TranslationRequest) error {
	translation := domain.NewTranslation(request)
	return s.repository.Save(*translation)
}

func (s *TranslationService) UpdateTranslation(id string, request domain.TranslationRequest) error {
	translation := s.repository.GetById(id)

	if translation == nil {
		return errors.New(fmt.Sprintf("Can not find translation by ID: %s", id))
	}

	translation.ApplyChanges(request)
	return s.repository.Save(*translation)
}

func (s *TranslationService) GetTranslations() []domain.Translation {
	return s.repository.Get()
}

func (s *TranslationService) GetById(id string) *domain.Translation {
	return s.repository.GetById(id)
}

func (s *TranslationService) DeleteById(id string) error {
	return s.repository.Delete(id)
}
