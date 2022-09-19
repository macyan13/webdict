package translation

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_CreateTranslation(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)

	mockCall := repository.On("Save", mock.Anything).Return(nil)
	err := service.CreateTranslation(Request{})
	assert.Nil(t, err)

	mockCall.Unset()

	repository.On("Save", mock.Anything).Return(errors.New("testError"))
	err = service.CreateTranslation(Request{})
	assert.Equal(t, "testError", err.Error())
}

func TestService_UpdateTranslation(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)

	id := "testId"
	request := Request{
		Transcription: "test",
		Translation:   "test",
		Text:          "test",
		Example:       "test",
	}

	mockGetByIdCall := repository.On("GetById", id).Times(1).Return(&Translation{})
	repository.On("Save", mock.MatchedBy(func(t Translation) bool { return t.Translation == "test" })).Times(1).Return(nil)
	err := service.UpdateTranslation(id, request)
	assert.Nil(t, err)

	mockGetByIdCall.Unset()

	repository.On("GetById", id).Times(1).Return(nil)
	err = service.UpdateTranslation(id, request)
	assert.Equal(t, fmt.Sprintf("Can not find translation by ID: %s", id), err.Error())
}

func TestService_GetTranslations(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)
	repository.On("Get").Times(1).Return([]Translation{})
	service.GetTranslations()
}

func TestService_GetById(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)
	id := "testId"

	repository.On("GetById", id).Times(1).Return(nil)
	translation := service.GetById(id)
	assert.Nil(t, translation)
}

func TestService_DeleteById(t *testing.T) {
	repository := MockRepository{}
	service := NewService(&repository)
	id := "testId"

	repository.On("Delete", id).Times(1).Return(nil)
	translation := service.DeleteById(id)
	assert.Nil(t, translation)

	repository.On("Delete", mock.Anything).Return(errors.New("testError"))
	err := service.DeleteById(id)
	assert.Equal(t, "testError", err.Error())
}
