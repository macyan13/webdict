// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/macyan13/webdict/backend/pkg/domain"
	mock "github.com/stretchr/testify/mock"
)

// TranslationRepository is an autogenerated mock type for the TranslationRepository type
type TranslationRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *TranslationRepository) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields:
func (_m *TranslationRepository) Get() []domain.Translation {
	ret := _m.Called()

	var r0 []domain.Translation
	if rf, ok := ret.Get(0).(func() []domain.Translation); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Translation)
		}
	}

	return r0
}

// GetById provides a mock function with given fields: id
func (_m *TranslationRepository) GetById(id string) *domain.Translation {
	ret := _m.Called(id)

	var r0 *domain.Translation
	if rf, ok := ret.Get(0).(func(string) *domain.Translation); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Translation)
		}
	}

	return r0
}

// Save provides a mock function with given fields: translation
func (_m *TranslationRepository) Save(translation domain.Translation) error {
	ret := _m.Called(translation)

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.Translation) error); ok {
		r0 = rf(translation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTranslationRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTranslationRepository creates a new instance of TranslationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTranslationRepository(t mockConstructorTestingTNewTranslationRepository) *TranslationRepository {
	mock := &TranslationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
