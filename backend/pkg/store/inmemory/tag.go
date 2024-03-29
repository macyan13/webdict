package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/query"
)

type TagRepo struct {
	storage map[string]*tag.Tag
}

func NewTagRepository() *TagRepo {
	return &TagRepo{
		storage: map[string]*tag.Tag{},
	}
}

func (r *TagRepo) Get(id, authorID string) (*tag.Tag, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		return t, nil
	}

	return nil, tag.ErrNotFound
}

func (r *TagRepo) Delete(id, authorID string) error {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		delete(r.storage, id)
		return nil
	}

	return fmt.Errorf("not found")
}

func (r *TagRepo) AllExist(ids []string, authorID string) (bool, error) {
	for _, id := range ids {
		if t, ok := r.storage[id]; !ok || t.AuthorID() != authorID {
			return false, nil
		}
	}

	return true, nil
}

func (r *TagRepo) Create(t *tag.Tag) error {
	r.storage[t.ID()] = t
	return nil
}

func (r *TagRepo) DeleteByAuthorID(authorID string) (int, error) {
	counter := 0
	for key, tg := range r.storage {
		if tg.AuthorID() == authorID {
			delete(r.storage, key)
			counter++
		}
	}

	return counter, nil
}

func (r *TagRepo) GetAllViews(authorID string) ([]query.TagView, error) {
	tags := make([]query.TagView, 0)
	for _, t := range r.storage {
		if t.AuthorID() != authorID {
			continue
		}

		tagData := t.ToMap()
		tags = append(tags, query.TagView{
			ID:   t.ID(),
			Name: tagData["name"].(string),
		})
	}

	return tags, nil
}

func (r *TagRepo) GetView(id, authorID string) (query.TagView, error) {
	t, ok := r.storage[id]

	if ok && t.AuthorID() == authorID {
		tagData := t.ToMap()
		return query.TagView{
			ID:   t.ID(),
			Name: tagData["name"].(string),
		}, nil
	}

	return query.TagView{}, fmt.Errorf("not found")
}

func (r *TagRepo) GetViews(ids []string, authorID string) ([]query.TagView, error) {
	views := make([]query.TagView, 0, len(ids))

	for _, id := range ids {
		for _, t := range r.storage {
			if t.AuthorID() == authorID && t.ID() == id {
				tagData := t.ToMap()
				views = append(views, query.TagView{
					ID:   t.ID(),
					Name: tagData["name"].(string),
				})
			}
		}
	}

	return views, nil
}

func (r *TagRepo) Update(t *tag.Tag) error {
	r.storage[t.ID()] = t
	return nil
}
