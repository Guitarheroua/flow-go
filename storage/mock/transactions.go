// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import crypto "github.com/dapperlabs/flow-go/crypto"
import flow "github.com/dapperlabs/flow-go/model/flow"
import mock "github.com/stretchr/testify/mock"

// Transactions is an autogenerated mock type for the Transactions type
type Transactions struct {
	mock.Mock
}

// ByHash provides a mock function with given fields: hash
func (_m *Transactions) ByHash(hash crypto.Hash) (*flow.Transaction, error) {
	ret := _m.Called(hash)

	var r0 *flow.Transaction
	if rf, ok := ret.Get(0).(func(crypto.Hash) *flow.Transaction); ok {
		r0 = rf(hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(crypto.Hash) error); ok {
		r1 = rf(hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: tx
func (_m *Transactions) Insert(tx *flow.Transaction) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Transaction) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
