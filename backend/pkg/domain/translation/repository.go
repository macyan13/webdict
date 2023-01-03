package translation

import "errors"

var ErrNotFound = errors.New("can not find translation in store")

// Repository defines domain translation repository methods
type Repository interface {
	Create(translation *Translation) error
	Update(translation *Translation) error
	Get(id, authorID string) (*Translation, error)   // Get provides translation by id and authorID, return ErrNotFound if record not exists
	ExistByText(text, authorID string) (bool, error) // ExistByText checks if translation with such text was already created
	ExistByTag(tagID, authorID string) (bool, error) // ExistByTag checks if at least one translation tagged with tagID exist
	Delete(id, authorID string) error
}
