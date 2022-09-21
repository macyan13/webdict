package translation

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type Repository interface {
	Save(translation Translation) error
	GetById(id string) *Translation
	Get() []Translation
	Delete(id string) error
}

type Service struct {
	repository    Repository
	tagRepository tag.Repository
}

func NewService(repository Repository, tagRepository tag.Repository) *Service {
	return &Service{
		repository:    repository,
		tagRepository: tagRepository,
	}
}

func (s *Service) CreateTranslation(request Request) error {
	data := data{}
	err := s.convertToData(request, &data)

	if err != nil {
		return err
	}

	translation := newTranslation(data)
	return s.repository.Save(*translation)
}

func (s *Service) UpdateTranslation(id string, request Request) error {
	translation := s.repository.GetById(id)

	if translation == nil {
		return errors.New(fmt.Sprintf("Can not find translation by ID: %s", id))
	}

	data := data{}
	err := s.convertToData(request, &data)

	if err != nil {
		return err
	}

	translation.applyChanges(data)
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

func (s *Service) convertToData(request Request, data *data) error {
	data.Request = request

	if len(request.TagIds) == 0 {
		return nil
	}

	existingTags := s.tagRepository.GetByIds(request.TagIds)

	if len(existingTags) != len(request.TagIds) {
		return errors.New("can not apply changes for translation tags, some passed tag are not found")
	}

	data.Tags = existingTags
	return nil
}
