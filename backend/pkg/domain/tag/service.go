package tag

import (
	"errors"
	"fmt"
)

type Repository interface {
	Save(tag Tag) error
	GetById(id string) *Tag
	GetByIds(ids []string) []*Tag // May be remove it
	Get() []Tag
	Delete(tag Tag) error
	AllExist(ids []string, AuthorId string) bool
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateTag(request Request) error {
	tag := NewTag(request)
	return s.repository.Save(*tag)
}

func (s *Service) UpdateTag(id string, request Request) error {
	tag := s.repository.GetById(id)

	if tag == nil {
		return errors.New(fmt.Sprintf("Can not find tag by ID: %s", id))
	}

	tag.ApplyChanges(request)
	return s.repository.Save(*tag)
}

func (s *Service) GetTags() []Tag {
	return s.repository.Get()
}

func (s *Service) GetById(id string) *Tag {
	return s.repository.GetById(id)
}

func (s *Service) DeleteById(id string) error {
	return s.repository.Delete(id)
}
