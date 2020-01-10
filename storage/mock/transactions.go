// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

// Transactions is an autogenerated mock type for the Transactions type
type Transactions struct {
	mock.Mock
}

// ByID provides a mock function with given fields: txID
func (_m *Transactions) ByID(txID flow.Identifier) (*flow.Transaction, error) {
	ret := _m.Called(txID)

	var r0 *flow.Transaction
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.Transaction); ok {
		r0 = rf(txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: txID
func (_m *Transactions) Remove(txID flow.Identifier) error {
	ret := _m.Called(txID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(txID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store provides a mock function with given fields: tx
func (_m *Transactions) Store(tx *flow.Transaction) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Transaction) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
