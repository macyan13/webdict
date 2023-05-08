package lang

import "errors"

var ErrNotFound = errors.New("can not find lang in store")

type Repository interface {
	Create(lang *Lang) error
	ExistByName(name, authorID string) (bool, error)
	Update(lang *Lang) error
	Get(id, authorID string) (*Lang, error)
	Delete(id, authorID string) error
}
