package tag

import "errors"

var ErrNotFound = errors.New("can not find tag in store")
var ErrTagAlreadyExists = errors.New("tag already exists")

type Repository interface {
	Create(tag *Tag) error                 // Create returns ErrTagAlreadyExists if record for pair tag-authorID already exists
	Update(tag *Tag) error                 // Update returns ErrTagAlreadyExists if record for pair tag-authorID already exists
	Get(id, authorID string) (*Tag, error) // Get provide tag by id and authorID, return ErrNotFound when tag not exist
	Delete(id, authorID string) error
	AllExist(ids []string, authorID string) (bool, error)
	DeleteByAuthorID(authorID string) (int, error)
}
