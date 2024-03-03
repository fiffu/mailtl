// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	backends "github.com/flashmob/go-guerrilla/backends"

	mail "github.com/flashmob/go-guerrilla/mail"

	mock "github.com/stretchr/testify/mock"
)

// Backend is an autogenerated mock type for the Backend type
type Backend struct {
	mock.Mock
}

// Initialize provides a mock function with given fields: backendConfig
func (_m *Backend) Initialize(backendConfig backends.BackendConfig) error {
	ret := _m.Called(backendConfig)

	if len(ret) == 0 {
		panic("no return value specified for Initialize")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(backends.BackendConfig) error); ok {
		r0 = rf(backendConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Backend) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SaveMail provides a mock function with given fields: ctx, e
func (_m *Backend) SaveMail(ctx context.Context, e *mail.Envelope) (bool, error) {
	ret := _m.Called(ctx, e)

	if len(ret) == 0 {
		panic("no return value specified for SaveMail")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *mail.Envelope) (bool, error)); ok {
		return rf(ctx, e)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *mail.Envelope) bool); ok {
		r0 = rf(ctx, e)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *mail.Envelope) error); ok {
		r1 = rf(ctx, e)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Shutdown provides a mock function with given fields:
func (_m *Backend) Shutdown() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Shutdown")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBackend creates a new instance of Backend. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBackend(t interface {
	mock.TestingT
	Cleanup(func())
}) *Backend {
	mock := &Backend{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
