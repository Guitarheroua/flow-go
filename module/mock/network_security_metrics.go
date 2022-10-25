// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import mock "github.com/stretchr/testify/mock"

// NetworkSecurityMetrics is an autogenerated mock type for the NetworkSecurityMetrics type
type NetworkSecurityMetrics struct {
	mock.Mock
}

// OnRateLimitedUnicastMessage provides a mock function with given fields: role, msgType, topic, reason
func (_m *NetworkSecurityMetrics) OnRateLimitedUnicastMessage(role string, msgType string, topic string, reason string) {
	_m.Called(role, msgType, topic, reason)
}

// OnUnauthorizedMessage provides a mock function with given fields: role, msgType, topic, offense
func (_m *NetworkSecurityMetrics) OnUnauthorizedMessage(role string, msgType string, topic string, offense string) {
	_m.Called(role, msgType, topic, offense)
}

type mockConstructorTestingTNewNetworkSecurityMetrics interface {
	mock.TestingT
	Cleanup(func())
}

// NewNetworkSecurityMetrics creates a new instance of NetworkSecurityMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNetworkSecurityMetrics(t mockConstructorTestingTNewNetworkSecurityMetrics) *NetworkSecurityMetrics {
	mock := &NetworkSecurityMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
