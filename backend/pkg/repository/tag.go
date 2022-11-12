package repository

import (
	"github.com/macyan13/webdict/backend/pkg/app/query"
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
	r.storage[tag.Id()] = tag
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

func (r *tagRepo) GetAll(authorId string) []query.Tag {
	tags := make([]query.Tag, 0)

	for _, t := range r.storage {
		if t.AuthorId() == authorId {
			tags = append(tags, query.Tag{
				Id:        t.Id(),
				Tag:       t.Tag(),
				CreatedAd: t.CreatedAt(),
			})
		}
	}

	return tags
}

func (r *tagRepo) GetTag(id, authorId string) *query.Tag {
	t, ok := r.storage[id]

	if ok && t.AuthorId() == authorId {
		return &query.Tag{
			Id:        t.Id(),
			Tag:       t.Tag(),
			CreatedAd: t.CreatedAt(),
		}
	}

	return nil
}

func (r *tagRepo) getTagsById(ids []string) []query.Tag {
	tags := make([]query.Tag, 0)

	for _, id := range ids {
		t, ok := r.storage[id]
		if ok {
			tags = append(tags, query.Tag{
				Id:        t.Id(),
				Tag:       t.Tag(),
				CreatedAd: t.CreatedAt(),
			})
		}
	}

	return tags
}

func (r *tagRepo) Delete(tag tag.Tag) error {
	delete(r.storage, tag.Id())
	return nil
}

func (r *tagRepo) AllExist(ids []string, AuthorId string) bool {
	for _, id := range ids {
		if _, ok := r.storage[id]; !ok {
			return false
		}
	}

	return true
}
