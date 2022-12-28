package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
)

type TagRepo struct {
	storage map[string]*tag.Tag
}

func NewTagRepository() *TagRepo {
	return &TagRepo{
		storage: map[string]*tag.Tag{},
	}
}

func (r *TagRepo) Get(id, authorId string) (tag.Tag, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorId {
		return *t, nil
	}

	return tag.Tag{}, tag.NotFoundErr
}

func (r *TagRepo) Delete(id, authorId string) error {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorId {
		delete(r.storage, id)
		return nil
	}

	return fmt.Errorf("not found")
}

func (r *TagRepo) AllExist(ids []string, authorId string) (bool, error) {
	for _, id := range ids {
		if t, ok := r.storage[id]; !ok || t.AuthorID() != authorId {
			return false, nil
		}
	}

	return true, nil
}

func (r *TagRepo) Create(t tag.Tag) error {
	r.storage[t.ID()] = &t
	return nil
}

func (r *TagRepo) GetAllViews(authorId string) ([]query.TagView, error) {
	tags := make([]query.TagView, 0)
	for _, t := range r.storage {
		if t.AuthorID() == authorId {
			tagData := t.ToMap()
			tags = append(tags, query.TagView{
				Id:  t.ID(),
				Tag: tagData["tag"].(string),
			})
		}
	}

	return tags, nil
}

func (r *TagRepo) GetView(id, authorId string) (query.TagView, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorId {
		tagData := t.ToMap()
		return query.TagView{
			Id:  t.ID(),
			Tag: tagData["tag"].(string),
		}, nil
	}

	return query.TagView{}, fmt.Errorf("not found")
}

func (r *TagRepo) GetViews(ids []string, authorId string) ([]query.TagView, error) {
	views := make([]query.TagView, 0, len(ids))

	for _, id := range ids {
		for _, t := range r.storage {
			if t.AuthorID() == authorId && t.ID() == id {
				tagData := t.ToMap()
				views = append(views, query.TagView{
					Id:  t.ID(),
					Tag: tagData["tag"].(string),
				})
			}
		}
	}

	return views, nil
}

func (r *TagRepo) Update(t tag.Tag) error {
	r.storage[t.ID()] = &t
	return nil
}
