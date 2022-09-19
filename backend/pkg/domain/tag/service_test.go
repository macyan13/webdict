package tag

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_CreateTag(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)

	mockCall := repository.On("Save", mock.Anything).Return(nil)
	err := service.CreateTag(Request{})
	assert.Nil(t, err)

	mockCall.Unset()

	repository.On("Save", mock.Anything).Return(errors.New("testError"))
	err = service.CreateTag(Request{})
	assert.Equal(t, "testError", err.Error())
}

func TestService_UpdateTag(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)

	id := "testId"
	request := Request{
		Tag: "test",
	}

	mockGetByIdCall := repository.On("GetById", id).Times(1).Return(&Tag{})
	repository.On("Save", mock.MatchedBy(func(t Tag) bool { return t.Tag == "test" })).Times(1).Return(nil)
	err := service.UpdateTag(id, request)
	assert.Nil(t, err)

	mockGetByIdCall.Unset()

	repository.On("GetById", id).Times(1).Return(nil)
	err = service.UpdateTag(id, request)
	assert.Equal(t, fmt.Sprintf("Can not find tag by ID: %s", id), err.Error())
}

func TestService_GetTag(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)
	repository.On("Get").Times(1).Return([]Tag{})
	service.GetTag()
}

func TestService_GetById(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)
	id := "testId"

	repository.On("GetById", id).Times(1).Return(nil)
	tag := service.GetById(id)
	assert.Nil(t, tag)
}

func TestService_DeleteById(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)
	id := "testId"

	repository.On("Delete", id).Times(1).Return(nil)
	tag := service.DeleteById(id)
	assert.Nil(t, tag)

	repository.On("Delete", mock.Anything).Return(errors.New("testError"))
	err := service.DeleteById(id)
	assert.Equal(t, "testError", err.Error())
}
