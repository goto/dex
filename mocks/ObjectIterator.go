// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	storage "cloud.google.com/go/storage"
	mock "github.com/stretchr/testify/mock"
)

// ObjectIterator is an autogenerated mock type for the ObjectIterator type
type ObjectIterator struct {
	mock.Mock
}

// Next provides a mock function with given fields:
func (_m *ObjectIterator) Next() (*storage.ObjectAttrs, error) {
	ret := _m.Called()

	var r0 *storage.ObjectAttrs
	var r1 error
	if rf, ok := ret.Get(0).(func() (*storage.ObjectAttrs, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *storage.ObjectAttrs); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.ObjectAttrs)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewObjectIterator creates a new instance of ObjectIterator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewObjectIterator(t interface {
	mock.TestingT
	Cleanup(func())
}) *ObjectIterator {
	mock := &ObjectIterator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
