// Code generated by mockery v2.13.0. DO NOT EDIT.

package mocks

import (
	hotstuff "github.com/onflow/flow-go/consensus/hotstuff"
	mock "github.com/stretchr/testify/mock"

	model "github.com/onflow/flow-go/consensus/hotstuff/model"
)

// VoteCollector is an autogenerated mock type for the VoteCollector type
type VoteCollector struct {
	mock.Mock
}

// AddVote provides a mock function with given fields: vote
func (_m *VoteCollector) AddVote(vote *model.Vote) error {
	ret := _m.Called(vote)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Vote) error); ok {
		r0 = rf(vote)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessBlock provides a mock function with given fields: block
func (_m *VoteCollector) ProcessBlock(block *model.Proposal) error {
	ret := _m.Called(block)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Proposal) error); ok {
		r0 = rf(block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterVoteConsumer provides a mock function with given fields: consumer
func (_m *VoteCollector) RegisterVoteConsumer(consumer hotstuff.VoteConsumer) {
	_m.Called(consumer)
}

// Status provides a mock function with given fields:
func (_m *VoteCollector) Status() hotstuff.VoteCollectorStatus {
	ret := _m.Called()

	var r0 hotstuff.VoteCollectorStatus
	if rf, ok := ret.Get(0).(func() hotstuff.VoteCollectorStatus); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(hotstuff.VoteCollectorStatus)
	}

	return r0
}

// View provides a mock function with given fields:
func (_m *VoteCollector) View() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

type NewVoteCollectorT interface {
	mock.TestingT
	Cleanup(func())
}

// NewVoteCollector creates a new instance of VoteCollector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVoteCollector(t NewVoteCollectorT) *VoteCollector {
	mock := &VoteCollector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
