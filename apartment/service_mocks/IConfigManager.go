// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IConfigManager is an autogenerated mock type for the IConfigManager type
type IConfigManager struct {
	mock.Mock
}

// GetAdminPassword provides a mock function with given fields:
func (_m *IConfigManager) GetAdminPassword() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAdminPassword")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetFromMail provides a mock function with given fields:
func (_m *IConfigManager) GetFromMail() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetFromMail")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetJwtKey provides a mock function with given fields:
func (_m *IConfigManager) GetJwtKey() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetJwtKey")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMailServer provides a mock function with given fields:
func (_m *IConfigManager) GetMailServer() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMailServer")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewIConfigManager creates a new instance of IConfigManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIConfigManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *IConfigManager {
	mock := &IConfigManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
