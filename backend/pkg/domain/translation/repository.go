package translation

import "errors"

var NotFoundErr = errors.New("can not find tag in storage")

// Repository defines domain translation repository methods
type Repository interface {
	Create(translation Translation) error
	Update(translation Translation) error
	Get(id, authorId string) (Translation, error) // Get provide translation by id and authorId, return NotFoundErr if record not exists
	Delete(id, authorId string) error
}
