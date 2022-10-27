package translation

import (
	"errors"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_CreateTranslation(t *testing.T) {
	translationRepository := MockRepository{}
	tagRepository := tag.MockRepository{}
	userRepository := user.MockRepository{}
	service := NewService(&translationRepository, &tagRepository, &userRepository)

	userRepository.On("Exist", mock.Anything).Times(2).Return(true)
	mockCall := translationRepository.On("Save", mock.Anything).Return(nil)
	err := service.CreateTranslation(Request{})
	assert.Nil(t, err)

	mockCall.Unset()
	translationRepository.On("Save", mock.Anything).Return(errors.New("testError"))

	err = service.CreateTranslation(Request{})
	assert.Equal(t, "testError", err.Error())
}

func TestService_CreateTranslationErrorOnMissedTagInDB(t *testing.T) {
	translationRepository := MockRepository{}
	tagRepository := tag.MockRepository{}
	userRepository := user.MockRepository{}
	service := NewService(&translationRepository, &tagRepository, &userRepository)

	tagRepository.On("GetByIds", mock.Anything).Return([]*tag.Tag{})
	userRepository.On("Exist", mock.Anything).Times(1).Return(true)

	err := service.CreateTranslation(Request{
		TagIds: []string{"notExistedTag"},
	})
	assert.Equal(t, "can not apply changes for translation tags, some passed tag are not found", err.Error())
}

func TestService_UpdateTranslation(t *testing.T) {
	translationRepository := MockRepository{}
	tagRepository := tag.MockRepository{}
	userRepository := user.MockRepository{}
	service := NewService(&translationRepository, &tagRepository, &userRepository)

	id := "testId"
	request := Request{
		Transcription: "test",
		Translation:   "test",
		Text:          "test",
		Example:       "test",
	}

	mockGetByIdCall := translationRepository.On("GetById", id).Times(1).Return(&Translation{})
	translationRepository.On("Save", mock.MatchedBy(func(t Translation) bool { return t.Translation == "test" })).Times(1).Return(nil)
	userRepository.On("Exist", mock.Anything).Times(1).Return(true)

	err := service.UpdateTranslation(id, request)
	assert.Nil(t, err)

	mockGetByIdCall.Unset()

	translationRepository.On("GetById", id).Times(1).Return(nil)
	err = service.UpdateTranslation(id, request)
	assert.Equal(t, fmt.Sprintf("can not find translation by ID: %s", id), err.Error())
}

func TestService_GetTranslations(t *testing.T) {
	translationRepository := MockRepository{}
	tagRepository := tag.MockRepository{}
	userRepository := user.MockRepository{}
	service := NewService(&translationRepository, &tagRepository, &userRepository)
	translationRepository.On("Get").Times(1).Return([]Translation{})
	service.GetTranslations()
}

func TestService_GetById(t *testing.T) {
	translationRepository := MockRepository{}
	tagRepository := tag.MockRepository{}
	userRepository := user.MockRepository{}
	service := NewService(&translationRepository, &tagRepository, &userRepository)
	id := "testId"

	translationRepository.On("GetById", id).Times(1).Return(nil)
	translation := service.GetById(id)
	assert.Nil(t, translation)
}

func TestService_DeleteById(t *testing.T) {
	translationRepository := MockRepository{}
	tagRepository := tag.MockRepository{}
	userRepository := user.MockRepository{}
	service := NewService(&translationRepository, &tagRepository, &userRepository)
	id := "testId"

	translationRepository.On("Delete", id).Times(1).Return(nil)
	translation := service.DeleteById(id)
	assert.Nil(t, translation)

	translationRepository.On("Delete", mock.Anything).Return(errors.New("testError"))
	err := service.DeleteById(id)
	assert.Equal(t, "testError", err.Error())
}
