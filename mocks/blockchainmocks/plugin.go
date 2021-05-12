// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package blockchainmocks

import (
	blockchain "github.com/kaleido-io/firefly/internal/blockchain"
	config "github.com/kaleido-io/firefly/internal/config"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *Plugin) Capabilities() *blockchain.Capabilities {
	ret := _m.Called()

	var r0 *blockchain.Capabilities
	if rf, ok := ret.Get(0).(func() *blockchain.Capabilities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*blockchain.Capabilities)
		}
	}

	return r0
}

// Init provides a mock function with given fields: ctx, _a1, events
func (_m *Plugin) Init(ctx context.Context, _a1 config.Config, events blockchain.Events) error {
	ret := _m.Called(ctx, _a1, events)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.Config, blockchain.Events) error); ok {
		r0 = rf(ctx, _a1, events)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitBroadcastBatch provides a mock function with given fields: ctx, identity, batch
func (_m *Plugin) SubmitBroadcastBatch(ctx context.Context, identity string, batch *blockchain.BroadcastBatch) (string, error) {
	ret := _m.Called(ctx, identity, batch)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, *blockchain.BroadcastBatch) string); ok {
		r0 = rf(ctx, identity, batch)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *blockchain.BroadcastBatch) error); ok {
		r1 = rf(ctx, identity, batch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyIdentitySyntax provides a mock function with given fields: ctx, identity
func (_m *Plugin) VerifyIdentitySyntax(ctx context.Context, identity string) (string, error) {
	ret := _m.Called(ctx, identity)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, identity)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, identity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
