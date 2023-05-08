// Code generated by mockery v2.21.4. DO NOT EDIT.

package mocks

import (
	model "github.com/onflow/flow-go/consensus/hotstuff/model"
	mock "github.com/stretchr/testify/mock"
)

// ProtocolViolationConsumer is an autogenerated mock type for the ProtocolViolationConsumer type
type ProtocolViolationConsumer struct {
	mock.Mock
}

// OnDoubleProposeDetected provides a mock function with given fields: _a0, _a1
func (_m *ProtocolViolationConsumer) OnDoubleProposeDetected(_a0 *model.Block, _a1 *model.Block) {
	_m.Called(_a0, _a1)
}

// OnDoubleTimeoutDetected provides a mock function with given fields: _a0, _a1
func (_m *ProtocolViolationConsumer) OnDoubleTimeoutDetected(_a0 *model.TimeoutObject, _a1 *model.TimeoutObject) {
	_m.Called(_a0, _a1)
}

// OnDoubleVotingDetected provides a mock function with given fields: _a0, _a1
func (_m *ProtocolViolationConsumer) OnDoubleVotingDetected(_a0 *model.Vote, _a1 *model.Vote) {
	_m.Called(_a0, _a1)
}

// OnInvalidBlockDetected provides a mock function with given fields: err
func (_m *ProtocolViolationConsumer) OnInvalidBlockDetected(err model.InvalidBlockError) {
	_m.Called(err)
}

// OnInvalidTimeoutDetected provides a mock function with given fields: err
func (_m *ProtocolViolationConsumer) OnInvalidTimeoutDetected(err model.InvalidTimeoutError) {
	_m.Called(err)
}

// OnInvalidVoteDetected provides a mock function with given fields: err
func (_m *ProtocolViolationConsumer) OnInvalidVoteDetected(err model.InvalidVoteError) {
	_m.Called(err)
}

// OnVoteForInvalidBlockDetected provides a mock function with given fields: vote, invalidProposal
func (_m *ProtocolViolationConsumer) OnVoteForInvalidBlockDetected(vote *model.Vote, invalidProposal *model.Proposal) {
	_m.Called(vote, invalidProposal)
}

type mockConstructorTestingTNewProtocolViolationConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewProtocolViolationConsumer creates a new instance of ProtocolViolationConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProtocolViolationConsumer(t mockConstructorTestingTNewProtocolViolationConsumer) *ProtocolViolationConsumer {
	mock := &ProtocolViolationConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
