package inmemory

import (
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/query"
)

type LangRepo struct {
	storage map[string]*lang.Lang
}

func NewLangRepository() *LangRepo {
	return &LangRepo{
		storage: map[string]*lang.Lang{},
	}
}

func (l LangRepo) Create(ln *lang.Lang) error {
	l.storage[ln.ID()] = ln
	return nil
}

func (l LangRepo) ExistByName(name, authorID string) (bool, error) {
	for _, ln := range l.storage {
		if ln.AuthorID() != authorID {
			continue
		}
		if ln.ToMap()["name"] == name {
			return true, nil
		}
	}
	return false, nil
}

func (l LangRepo) Update(ln *lang.Lang) error {
	l.storage[ln.ID()] = ln
	return nil
}

func (l LangRepo) Get(id, authorID string) (*lang.Lang, error) {
	ln, ok := l.storage[id]

	if ok && ln.AuthorID() == authorID {
		return ln, nil
	}

	return nil, lang.ErrNotFound
}

func (l LangRepo) Delete(id, authorID string) error {
	ln, ok := l.storage[id]

	if ok && ln.AuthorID() == authorID {
		delete(l.storage, id)
		return nil
	}

	return fmt.Errorf("not found")
}

func (l LangRepo) Exist(id, authorID string) (bool, error) {
	ln, ok := l.storage[id]

	if ok && ln.AuthorID() == authorID {
		return true, nil
	}

	return false, nil
}

func (l LangRepo) GetAllViews(authorID string) ([]query.LangView, error) {
	langs := make([]query.LangView, 0)
	for _, ln := range l.storage {
		if ln.AuthorID() != authorID {
			continue
		}

		langData := ln.ToMap()
		langs = append(langs, query.LangView{
			ID:   ln.ID(),
			Name: langData["name"].(string),
		})
	}

	return langs, nil
}

func (l LangRepo) GetView(id, authorID string) (query.LangView, error) {
	ln, ok := l.storage[id]

	if ok && ln.AuthorID() == authorID {
		langData := ln.ToMap()
		return query.LangView{
			ID:   ln.ID(),
			Name: langData["name"].(string),
		}, nil
	}

	return query.LangView{}, fmt.Errorf("not found")
}
