package tag

import "errors"

var NotFoundErr = errors.New("can not find tag in storage")

type Repository interface {
	Create(tag Tag) error
	Update(tag Tag) error
	Get(id, authorId string) (Tag, error) // Get provide tag by id and authorID, return NotFoundErr when tag not exist
	Delete(id, authorId string) error
	AllExist(ids []string, authorId string) (bool, error)
}
