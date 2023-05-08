package lang

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Lang struct {
	id       string
	name     string
	authorID string
}

func NewLang(name, authorID string) (*Lang, error) {
	ln := Lang{
		id:       uuid.New().String(),
		name:     name,
		authorID: authorID,
	}

	if err := ln.validate(); err != nil {
		return nil, err
	}

	return &ln, nil
}

func (l *Lang) ID() string {
	return l.id
}

func (l *Lang) AuthorID() string {
	return l.authorID
}

func (l *Lang) Name() string {
	return l.name
}

func (l *Lang) ApplyChanges(name string) error {
	updated := *l
	updated.applyChanges(name)

	if err := updated.validate(); err != nil {
		return err
	}

	l.applyChanges(name)
	return nil
}

func (l *Lang) applyChanges(name string) {
	l.name = name
}

func (l *Lang) validate() error {
	var err error
	if l.name == "" {
		err = errors.Join(errors.New("name can not be empty"), err)
	}

	if l.authorID == "" {
		err = errors.Join(fmt.Errorf("authorID can not be empty"), err)
	}

	return err
}

func (l *Lang) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       l.id,
		"name":     l.name,
		"authorID": l.authorID,
	}
}

func UnmarshalFromDB(
	id string,
	name string,
	authorID string,
) *Lang {
	return &Lang{
		id:       id,
		name:     name,
		authorID: authorID,
	}
}
