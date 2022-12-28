package tag

import "errors"

var ErrNotFound = errors.New("can not find tag in storage")

type Repository interface {
	Create(tag Tag) error
	Update(tag Tag) error
	Get(id, authorID string) (Tag, error) // Get provide tag by id and authorID, return ErrNotFound when tag not exist
	Delete(id, authorID string) error
	AllExist(ids []string, authorID string) (bool, error)
}
