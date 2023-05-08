// Code generated by mockery v2.23.1. DO NOT EDIT.

package query

import mock "github.com/stretchr/testify/mock"

// mockery --name=LangViewRepository --filename=tlang_view_repository_mock.go --output=./ --structname=MockLangViewRepository --inpackage
// MockLangViewRepository is an autogenerated mock type for the LangViewRepository type
type MockLangViewRepository struct {
	mock.Mock
}

// GetAllViews provides a mock function with given fields: authorID
func (_m *MockLangViewRepository) GetAllViews(authorID string) ([]LangView, error) {
	ret := _m.Called(authorID)

	var r0 []LangView
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]LangView, error)); ok {
		return rf(authorID)
	}
	if rf, ok := ret.Get(0).(func(string) []LangView); ok {
		r0 = rf(authorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]LangView)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(authorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetView provides a mock function with given fields: id, authorID
func (_m *MockLangViewRepository) GetView(id string, authorID string) (LangView, error) {
	ret := _m.Called(id, authorID)

	var r0 LangView
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (LangView, error)); ok {
		return rf(id, authorID)
	}
	if rf, ok := ret.Get(0).(func(string, string) LangView); ok {
		r0 = rf(id, authorID)
	} else {
		r0 = ret.Get(0).(LangView)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, authorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockLangViewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockLangViewRepository creates a new instance of MockLangViewRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockLangViewRepository(t mockConstructorTestingTNewMockLangViewRepository) *MockLangViewRepository {
	mock := &MockLangViewRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
