// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	backends "github.com/flashmob/go-guerrilla/backends"

	mail "github.com/flashmob/go-guerrilla/mail"

	mock "github.com/stretchr/testify/mock"
)

// FixtureBackend is an autogenerated mock type for the FixtureBackend type
type FixtureBackend struct {
	mock.Mock
}

// Process provides a mock function with given fields: _a0, _a1
func (_m *FixtureBackend) Process(_a0 *mail.Envelope, _a1 backends.SelectTask) (backends.Result, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Process")
	}

	var r0 backends.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(*mail.Envelope, backends.SelectTask) (backends.Result, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(*mail.Envelope, backends.SelectTask) backends.Result); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(backends.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(*mail.Envelope, backends.SelectTask) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFixtureBackend creates a new instance of FixtureBackend. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFixtureBackend(t interface {
	mock.TestingT
	Cleanup(func())
}) *FixtureBackend {
	mock := &FixtureBackend{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}