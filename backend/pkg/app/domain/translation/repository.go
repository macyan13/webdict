package translation

import "errors"

var ErrNotFound = errors.New("can not find translation in store")

// Repository defines domain translation repository methods
type Repository interface {
	Create(translation *Translation) error
	Update(translation *Translation) error                       // Update saves the updated translation entity to store
	Get(id, authorID string) (*Translation, error)               // Get provides translation by id and authorID, return ErrNotFound if record not exists
	ExistBySource(source, authorID, langID string) (bool, error) // ExistByText checks if translation with such source was already created
	ExistByTag(tagID, authorID string) (bool, error)             // ExistByTag checks if at least one translation tagged with tagID exist
	ExistByLang(langID, authorID string) (bool, error)           // ExistByLang checks if at least one translation created with the passed language
	Delete(id, authorID string) error
}
