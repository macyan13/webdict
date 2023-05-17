package lang

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLang_ApplyChanges(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Error on validation",
			args{name: ""},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name can not be empty", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			args{name: "de"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			l := &Lang{
				name:     "en",
				authorID: "authorID",
			}
			if tt.wantErr(t1, l.ApplyChanges(tt.args.name), fmt.Sprintf("ApplyChanges(%v)", tt.args.name)) {
				assert.Equal(t1, "en", l.name)
			}
		})
	}
}

func TestNewLang(t *testing.T) {
	type args struct {
		name     string
		authorID string
	}
	tests := []struct {
		name    string
		args    args
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Error on validation",
			args{
				name:     "",
				authorID: "test",
			},
			nil,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name can not be empty", err.Error(), i)
				return false
			},
		},
		{
			"Positive case",
			args{
				name:     "EN",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, ln interface{}, i ...interface{}) bool {
				result := ln.(*Lang)
				assert.Equal(t, "EN", result.Name(), i)
				assert.Equal(t, "testAuthor", result.AuthorID(), i)
				return false
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLang(tt.args.name, tt.args.authorID)
			if !tt.wantErr(t, err, fmt.Sprintf("NewLang(%v, %v)", tt.args.name, tt.args.authorID)) {
				return
			}
			tt.want(t, got, "NewLang(%v, %v)", tt.args.name, tt.args.authorID)
		})
	}
}

func TestLang_validate(t *testing.T) {
	type fields struct {
		name     string
		authorID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"AuthorID is empty",
			fields{
				name:     "en",
				authorID: "",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "authorID can not be empty", err.Error(), i)
				return true
			},
		},
		{
			"Name is empty",
			fields{
				name:     "",
				authorID: "authorID",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name can not be empty", err.Error(), i)
				return true
			},
		},
		{
			"Multiple errors",
			fields{
				name:     "",
				authorID: "",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "authorID can not be empty"), i)
				assert.True(t, strings.Contains(err.Error(), "name can not be empty"), i)
				return true
			},
		},
		{
			"Name is empty",
			fields{
				name:     "",
				authorID: "authorID",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name can not be empty", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			fields{
				name:     "en",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln := &Lang{
				name:     tt.fields.name,
				authorID: tt.fields.authorID,
			}
			tt.wantErr(t, ln.validate(), "validate()")
		})
	}
}

func TestUnmarshalFromDB(t *testing.T) {
	ln := Lang{
		id:       "testId",
		name:     "testLang",
		authorID: "testAuthor",
	}

	assert.Equal(t, &ln, UnmarshalFromDB(ln.id, ln.name, ln.authorID))
}
