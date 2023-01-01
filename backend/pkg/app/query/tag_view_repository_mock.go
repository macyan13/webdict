// Code generated by mockery v2.14.0. DO NOT EDIT.

package query

import mock "github.com/stretchr/testify/mock"

//  mockery --name=TagViewRepository --filename=tag_view_repository_mock.go --output=./ --structname=MockTagViewRepository --inpackage
// MockTagViewRepository is an autogenerated mock type for the TagViewRepository type
type MockTagViewRepository struct {
	mock.Mock
}

// GetAllViews provides a mock function with given fields: authorId
func (_m *MockTagViewRepository) GetAllViews(authorId string) ([]TagView, error) {
	ret := _m.Called(authorId)

	var r0 []TagView
	if rf, ok := ret.Get(0).(func(string) []TagView); ok {
		r0 = rf(authorId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]TagView)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(authorId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetView provides a mock function with given fields: id, authorId
func (_m *MockTagViewRepository) GetView(id string, authorId string) (TagView, error) {
	ret := _m.Called(id, authorId)

	var r0 TagView
	if rf, ok := ret.Get(0).(func(string, string) TagView); ok {
		r0 = rf(id, authorId)
	} else {
		r0 = ret.Get(0).(TagView)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, authorId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetViews provides a mock function with given fields: ids, authorId
func (_m *MockTagViewRepository) GetViews(ids []string, authorId string) ([]TagView, error) {
	ret := _m.Called(ids, authorId)

	var r0 []TagView
	if rf, ok := ret.Get(0).(func([]string, string) []TagView); ok {
		r0 = rf(ids, authorId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]TagView)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, string) error); ok {
		r1 = rf(ids, authorId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockTagViewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockTagViewRepository creates a new instance of MockTagViewRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTagViewRepository(t mockConstructorTestingTNewMockTagViewRepository) *MockTagViewRepository {
	mock := &MockTagViewRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}