// Code generated by mockery v2.14.0. DO NOT EDIT.

package tag

import (
	mock "github.com/stretchr/testify/mock"
)

// MockRepository is an autogenerated mock type for the MockRepository type
type MockRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *MockRepository) Delete(id string) error {
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
func (_m *MockRepository) Get() []Tag {
	ret := _m.Called()

	var r0 []Tag
	if rf, ok := ret.Get(0).(func() []Tag); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Tag)
		}
	}

	return r0
}

// GetById provides a mock function with given fields: id
func (_m *MockRepository) GetById(id string) *Tag {
	ret := _m.Called(id)

	var r0 *Tag
	if rf, ok := ret.Get(0).(func(string) *Tag); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Tag)
		}
	}

	return r0
}

// GetByIds provides a mock function with given fields: ids
func (_m *MockRepository) GetByIds(ids []string) []*Tag {
	ret := _m.Called(ids)

	var r0 []*Tag
	if rf, ok := ret.Get(0).(func([]string) []*Tag); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Tag)
		}
	}

	return r0
}

// Save provides a mock function with given fields: _a0
func (_m *MockRepository) Save(_a0 Tag) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(Tag) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
