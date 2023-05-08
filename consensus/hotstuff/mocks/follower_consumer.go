// Code generated by mockery v2.21.4. DO NOT EDIT.

package mocks

import (
	model "github.com/onflow/flow-go/consensus/hotstuff/model"
	mock "github.com/stretchr/testify/mock"
)

// FollowerConsumer is an autogenerated mock type for the FollowerConsumer type
type FollowerConsumer struct {
	mock.Mock
}

// OnBlockIncorporated provides a mock function with given fields: _a0
func (_m *FollowerConsumer) OnBlockIncorporated(_a0 *model.Block) {
	_m.Called(_a0)
}

// OnDoubleProposeDetected provides a mock function with given fields: _a0, _a1
func (_m *FollowerConsumer) OnDoubleProposeDetected(_a0 *model.Block, _a1 *model.Block) {
	_m.Called(_a0, _a1)
}

// OnDoubleTimeoutDetected provides a mock function with given fields: _a0, _a1
func (_m *FollowerConsumer) OnDoubleTimeoutDetected(_a0 *model.TimeoutObject, _a1 *model.TimeoutObject) {
	_m.Called(_a0, _a1)
}

// OnDoubleVotingDetected provides a mock function with given fields: _a0, _a1
func (_m *FollowerConsumer) OnDoubleVotingDetected(_a0 *model.Vote, _a1 *model.Vote) {
	_m.Called(_a0, _a1)
}

// OnFinalizedBlock provides a mock function with given fields: _a0
func (_m *FollowerConsumer) OnFinalizedBlock(_a0 *model.Block) {
	_m.Called(_a0)
}

// OnInvalidBlockDetected provides a mock function with given fields: err
func (_m *FollowerConsumer) OnInvalidBlockDetected(err model.InvalidBlockError) {
	_m.Called(err)
}

// OnInvalidTimeoutDetected provides a mock function with given fields: err
func (_m *FollowerConsumer) OnInvalidTimeoutDetected(err model.InvalidTimeoutError) {
	_m.Called(err)
}

// OnInvalidVoteDetected provides a mock function with given fields: err
func (_m *FollowerConsumer) OnInvalidVoteDetected(err model.InvalidVoteError) {
	_m.Called(err)
}

// OnVoteForInvalidBlockDetected provides a mock function with given fields: vote, invalidProposal
func (_m *FollowerConsumer) OnVoteForInvalidBlockDetected(vote *model.Vote, invalidProposal *model.Proposal) {
	_m.Called(vote, invalidProposal)
}

type mockConstructorTestingTNewFollowerConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewFollowerConsumer creates a new instance of FollowerConsumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFollowerConsumer(t mockConstructorTestingTNewFollowerConsumer) *FollowerConsumer {
	mock := &FollowerConsumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
