package translation

import "errors"

var ErrNotFound = errors.New("can not find translation in store")
var ErrSourceAlreadyExists = errors.New("translation with such source already exists")

// Repository defines domain translation repository methods
type Repository interface {
	Create(translation *Translation) error             // Create returns ErrSourceAlreadyExists if records with values for source-langId-authorID already exists
	Update(translation *Translation) error             // Update saves the updated translation entity to store, returns ErrSourceAlreadyExists if records with values for source-langId-authorID already exists
	Get(id, authorID string) (*Translation, error)     // Get provides translation by id and authorID, return ErrNotFound if record not exists
	ExistByTag(tagID, authorID string) (bool, error)   // ExistByTag checks if at least one translation tagged with tagID exist
	ExistByLang(langID, authorID string) (bool, error) // ExistByLang checks if at least one translation created with the passed language
	Delete(id, authorID string) error
}
