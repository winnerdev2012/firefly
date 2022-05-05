// Code generated by mockery v1.0.0. DO NOT EDIT.

package identitymocks

import (
	context "context"

	config "github.com/hyperledger/firefly-common/pkg/config"

	identity "github.com/hyperledger/firefly/pkg/identity"

	mock "github.com/stretchr/testify/mock"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *Plugin) Capabilities() *identity.Capabilities {
	ret := _m.Called()

	var r0 *identity.Capabilities
	if rf, ok := ret.Get(0).(func() *identity.Capabilities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*identity.Capabilities)
		}
	}

	return r0
}

// Init provides a mock function with given fields: ctx, prefix, callbacks
func (_m *Plugin) Init(ctx context.Context, prefix config.Prefix, callbacks identity.Callbacks) error {
	ret := _m.Called(ctx, prefix, callbacks)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.Prefix, identity.Callbacks) error); ok {
		r0 = rf(ctx, prefix, callbacks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitPrefix provides a mock function with given fields: prefix
func (_m *Plugin) InitPrefix(prefix config.Prefix) {
	_m.Called(prefix)
}

// Name provides a mock function with given fields:
func (_m *Plugin) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *Plugin) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
