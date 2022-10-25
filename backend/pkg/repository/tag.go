package repository

import (
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type tagRepo struct {
	storage map[string]tag.Tag
}

func NewTagRepository() *tagRepo {
	return &tagRepo{
		storage: map[string]tag.Tag{},
	}
}

func (r *tagRepo) Save(tag tag.Tag) error {
	r.storage[tag.Id] = tag
	return nil
}

func (r *tagRepo) Get() []tag.Tag {
	result := make([]tag.Tag, 0, len(r.storage))

	for _, tag := range r.storage {
		result = append(result, tag)
	}

	return result
}

func (r *tagRepo) GetById(id string) *tag.Tag {
	tag, ok := r.storage[id]

	if ok {
		return &tag
	}

	return nil
}

func (r *tagRepo) Delete(id string) error {
	delete(r.storage, id)
	return nil
}

func (r *tagRepo) GetByIds(ids []string) []*tag.Tag {
	var tags []*tag.Tag
	for _, id := range ids {
		el, ok := r.storage[id]

		if ok {
			tags = append(tags, &el)
		}
	}

	return tags
}
