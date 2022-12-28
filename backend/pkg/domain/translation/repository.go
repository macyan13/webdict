package translation

import "errors"

var ErrNotFound = errors.New("can not find tag in storage")

// Repository defines domain translation repository methods
type Repository interface {
	Create(translation Translation) error
	Update(translation Translation) error
	Get(id, authorID string) (Translation, error) // Get provide translation by id and authorID, return ErrNotFound if record not exists
	Delete(id, authorID string) error
}
