package tag

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTag_ApplyChanges(t *testing.T) {
	type fields struct {
		tag      string
		authorID string
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Error on validation",
			fields{
				tag:      "tat",
				authorID: "test",
			},
			args{tag: "t"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name length should be at least 2 symbols, 1 passed (t)", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			fields{
				tag:      "testTag",
				authorID: "testAuthor",
			},
			args{tag: "name"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			t := &Tag{
				name:     tt.fields.tag,
				authorID: tt.fields.authorID,
			}
			if tt.wantErr(t1, t.ApplyChanges(tt.args.tag), fmt.Sprintf("ApplyChanges(%v)", tt.args.tag)) {
				assert.Equal(t1, tt.fields.tag, t.name)
			}
		})
	}
}

func TestUnmarshalFromDB(t *testing.T) {
	tag := Tag{
		id:       "testId",
		name:     "testTag",
		authorID: "testAuthor",
	}

	assert.Equal(t, &tag, UnmarshalFromDB(tag.id, tag.name, tag.authorID))
}

func TestNewTag(t *testing.T) {
	type args struct {
		tag      string
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
				tag:      "t",
				authorID: "test",
			},
			nil,
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name length should be at least 2 symbols, 1 passed (t)", err.Error(), i)
				return false
			},
		},
		{
			"Positive case",
			args{
				tag:      "testTag",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, tg interface{}, i ...interface{}) bool {
				result := tg.(*Tag)
				assert.Equal(t, "testTag", result.name, i)
				assert.Equal(t, "testAuthor", result.authorID, i)
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
			got, err := NewTag(tt.args.tag, tt.args.authorID)
			if !tt.wantErr(t, err, fmt.Sprintf("NewTag(%v, %v)", tt.args.tag, tt.args.authorID)) {
				return
			}
			tt.want(t, got, "NewTag(%v, %v)", tt.args.tag, tt.args.authorID)
		})
	}
}

func TestTag_validate(t *testing.T) {
	type fields struct {
		tag      string
		authorID string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"Name is too short",
			fields{
				tag:      "t",
				authorID: "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "name length should be at least 2 symbols, 1 passed (t)", err.Error(), i)
				return true
			},
		},
		{
			"Name is too long",
			fields{
				tag:      string(make([]rune, 31)),
				authorID: "test",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.True(t, strings.Contains(err.Error(), "name max length is 30 symbols, 31 passed"), i)
				return true
			},
		},
		{
			"AuthorID is empty",
			fields{
				tag:      "tag",
				authorID: "",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "authorID can not be empty", err.Error(), i)
				return true
			},
		},
		{
			"Positive case",
			fields{
				tag:      "testTag",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			t := &Tag{
				name:     tt.fields.tag,
				authorID: tt.fields.authorID,
			}
			tt.wantErr(t1, t.validate(), "validate()")
		})
	}
}
