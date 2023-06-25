package lang

import "errors"

var ErrNotFound = errors.New("can not find lang in store")
var ErrLangAlreadyExists = errors.New("lang already exists")

type Repository interface {
	Create(lang *Lang) error // Create returns ErrLangAlreadyExists if record for pair name-authorID already exists
	Exist(id, authorID string) (bool, error)
	Update(lang *Lang) error // Update returns ErrLangAlreadyExists if record for pair name-authorID already exists
	Get(id, authorID string) (*Lang, error)
	Delete(id, authorID string) error
	DeleteByAuthorID(authorID string) (int, error)
}
