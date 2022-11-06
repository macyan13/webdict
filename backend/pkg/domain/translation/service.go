package translation

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
)

type Repository interface {
	Save(translation Translation) error
	GetById(id string) *Translation
	Get() []Translation
	Delete(translation Translation) error
}

type Service struct {
	repository     Repository
	tagRepository  tag.Repository
	userRepository user.Repository
}

func NewService(repository Repository, tagRepository tag.Repository, userRepository user.Repository) *Service {
	return &Service{
		repository:     repository,
		tagRepository:  tagRepository,
		userRepository: userRepository,
	}
}

func (s *Service) CreateTranslation(request Request) error {
	if err := s.validateRequest(request); err != nil {
		return err
	}

	data := Data{}

	if err := s.convertToData(request, &data); err != nil {
		return err
	}

	translation := NewTranslation(data)
	return s.repository.Save(*translation)
}

func (s *Service) UpdateTranslation(id string, request Request) error {
	translation := s.repository.GetById(id)

	if translation == nil {
		return fmt.Errorf("can not find translation by ID: %s", id)
	}

	data := Data{}

	if err := s.convertToData(request, &data); err != nil {
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

func (s *Service) validateRequest(request Request) error {
	if !s.userRepository.Exist(request.AuthorId) {
		return fmt.Errorf("can not find user by passed AuthorId: %s", request.AuthorId)
	}
	return nil
}

func (s *Service) convertToData(request Request, data *Data) error {
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
