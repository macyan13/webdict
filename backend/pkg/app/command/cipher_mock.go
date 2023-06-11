// Code generated by mockery v2.23.1. DO NOT EDIT.

package command

import mock "github.com/stretchr/testify/mock"

// mockery --name=Cipher --filename=cipher_mock.go --output=./ --structname=MockCipher --inpackage
// MockCipher is an autogenerated mock type for the Cipher type
type MockCipher struct {
	mock.Mock
}

// ComparePasswords provides a mock function with given fields: hashedPwd, plainPwd
func (_m *MockCipher) ComparePasswords(hashedPwd string, plainPwd string) bool {
	ret := _m.Called(hashedPwd, plainPwd)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(hashedPwd, plainPwd)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GenerateHash provides a mock function with given fields: pwd
func (_m *MockCipher) GenerateHash(pwd string) (string, error) {
	ret := _m.Called(pwd)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(pwd)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(pwd)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(pwd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockCipher interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockCipher creates a new instance of MockCipher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCipher(t mockConstructorTestingTNewMockCipher) *MockCipher {
	mock := &MockCipher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
