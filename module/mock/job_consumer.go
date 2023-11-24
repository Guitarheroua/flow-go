// Code generated by mockery v2.21.4. DO NOT EDIT.

package mock

import (
	module "github.com/onflow/flow-go/module"
	mock "github.com/stretchr/testify/mock"
)

// JobConsumer is an autogenerated mock type for the JobConsumer type
type JobConsumer struct {
	mock.Mock
}

// Check provides a mock function with given fields:
func (_m *JobConsumer) Check() {
	_m.Called()
}

// LastProcessedIndex provides a mock function with given fields:
func (_m *JobConsumer) LastProcessedIndex() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// NotifyJobIsDone provides a mock function with given fields: _a0
func (_m *JobConsumer) NotifyJobIsDone(_a0 module.JobID) uint64 {
	ret := _m.Called(_a0)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(module.JobID) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *JobConsumer) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *JobConsumer) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *JobConsumer) Stop() {
	_m.Called()
}

type mockConstructorTestingTNewJobConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewJobConsumer creates a new instance of JobConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewJobConsumer(t mockConstructorTestingTNewJobConsumer) *JobConsumer {
	mock := &JobConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
