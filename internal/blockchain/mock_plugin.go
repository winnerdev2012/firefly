// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchain

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockPlugin is an autogenerated mock type for the Plugin type
type MockPlugin struct {
	mock.Mock
}

// ConfigInterface provides a mock function with given fields:
func (_m *MockPlugin) ConfigInterface() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Init provides a mock function with given fields: ctx, config, events
func (_m *MockPlugin) Init(ctx context.Context, config interface{}, events Events) (*Capabilities, error) {
	ret := _m.Called(ctx, config, events)

	var r0 *Capabilities
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, Events) *Capabilities); ok {
		r0 = rf(ctx, config, events)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Capabilities)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, Events) error); ok {
		r1 = rf(ctx, config, events)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitBroadcastBatch provides a mock function with given fields: ctx, identity, broadcast
func (_m *MockPlugin) SubmitBroadcastBatch(ctx context.Context, identity string, broadcast *BroadcastBatch) (string, error) {
	ret := _m.Called(ctx, identity, broadcast)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, *BroadcastBatch) string); ok {
		r0 = rf(ctx, identity, broadcast)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *BroadcastBatch) error); ok {
		r1 = rf(ctx, identity, broadcast)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
