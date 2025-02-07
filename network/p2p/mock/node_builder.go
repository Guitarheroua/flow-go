// Code generated by mockery v2.43.2. DO NOT EDIT.

package mockp2p

import (
	context "context"

	connmgr "github.com/libp2p/go-libp2p/core/connmgr"

	host "github.com/libp2p/go-libp2p/core/host"

	madns "github.com/multiformats/go-multiaddr-dns"

	mock "github.com/stretchr/testify/mock"

	network "github.com/libp2p/go-libp2p/core/network"

	p2p "github.com/onflow/flow-go/network/p2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	routing "github.com/libp2p/go-libp2p/core/routing"
)

// NodeBuilder is an autogenerated mock type for the NodeBuilder type
type NodeBuilder struct {
	mock.Mock
}

// Build provides a mock function with given fields:
func (_m *NodeBuilder) Build() (p2p.LibP2PNode, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Build")
	}

	var r0 p2p.LibP2PNode
	var r1 error
	if rf, ok := ret.Get(0).(func() (p2p.LibP2PNode, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() p2p.LibP2PNode); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.LibP2PNode)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OverrideDefaultRpcInspectorFactory provides a mock function with given fields: _a0
func (_m *NodeBuilder) OverrideDefaultRpcInspectorFactory(_a0 p2p.GossipSubRpcInspectorFactoryFunc) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for OverrideDefaultRpcInspectorFactory")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(p2p.GossipSubRpcInspectorFactoryFunc) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// OverrideDefaultValidateQueueSize provides a mock function with given fields: _a0
func (_m *NodeBuilder) OverrideDefaultValidateQueueSize(_a0 int) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for OverrideDefaultValidateQueueSize")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(int) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// OverrideGossipSubFactory provides a mock function with given fields: _a0, _a1
func (_m *NodeBuilder) OverrideGossipSubFactory(_a0 p2p.GossipSubFactoryFunc, _a1 p2p.GossipSubAdapterConfigFunc) p2p.NodeBuilder {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for OverrideGossipSubFactory")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(p2p.GossipSubFactoryFunc, p2p.GossipSubAdapterConfigFunc) p2p.NodeBuilder); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// OverrideGossipSubScoringConfig provides a mock function with given fields: _a0
func (_m *NodeBuilder) OverrideGossipSubScoringConfig(_a0 *p2p.PeerScoringConfigOverride) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for OverrideGossipSubScoringConfig")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(*p2p.PeerScoringConfigOverride) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// OverrideNodeConstructor provides a mock function with given fields: _a0
func (_m *NodeBuilder) OverrideNodeConstructor(_a0 p2p.NodeConstructor) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for OverrideNodeConstructor")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(p2p.NodeConstructor) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// SetBasicResolver provides a mock function with given fields: _a0
func (_m *NodeBuilder) SetBasicResolver(_a0 madns.BasicResolver) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetBasicResolver")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(madns.BasicResolver) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// SetConnectionGater provides a mock function with given fields: _a0
func (_m *NodeBuilder) SetConnectionGater(_a0 p2p.ConnectionGater) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetConnectionGater")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(p2p.ConnectionGater) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// SetConnectionManager provides a mock function with given fields: _a0
func (_m *NodeBuilder) SetConnectionManager(_a0 connmgr.ConnManager) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetConnectionManager")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(connmgr.ConnManager) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// SetResourceManager provides a mock function with given fields: _a0
func (_m *NodeBuilder) SetResourceManager(_a0 network.ResourceManager) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetResourceManager")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(network.ResourceManager) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// SetRoutingSystem provides a mock function with given fields: _a0
func (_m *NodeBuilder) SetRoutingSystem(_a0 func(context.Context, host.Host) (routing.Routing, error)) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetRoutingSystem")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(func(context.Context, host.Host) (routing.Routing, error)) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// SetSubscriptionFilter provides a mock function with given fields: _a0
func (_m *NodeBuilder) SetSubscriptionFilter(_a0 pubsub.SubscriptionFilter) p2p.NodeBuilder {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SetSubscriptionFilter")
	}

	var r0 p2p.NodeBuilder
	if rf, ok := ret.Get(0).(func(pubsub.SubscriptionFilter) p2p.NodeBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.NodeBuilder)
		}
	}

	return r0
}

// NewNodeBuilder creates a new instance of NodeBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNodeBuilder(t interface {
	mock.TestingT
	Cleanup(func())
}) *NodeBuilder {
	mock := &NodeBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
