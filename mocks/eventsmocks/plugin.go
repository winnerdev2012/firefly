// Code generated by mockery v1.0.0. DO NOT EDIT.

package eventsmocks

import (
	context "context"

	config "github.com/hyperledger/firefly-common/pkg/config"

	core "github.com/hyperledger/firefly/pkg/core"

	events "github.com/hyperledger/firefly/pkg/events"

	mock "github.com/stretchr/testify/mock"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *Plugin) Capabilities() *events.Capabilities {
	ret := _m.Called()

	var r0 *events.Capabilities
	if rf, ok := ret.Get(0).(func() *events.Capabilities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*events.Capabilities)
		}
	}

	return r0
}

// DeliveryRequest provides a mock function with given fields: connID, sub, event, data
func (_m *Plugin) DeliveryRequest(connID string, sub *core.Subscription, event *core.EventDelivery, data core.DataArray) error {
	ret := _m.Called(connID, sub, event, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *core.Subscription, *core.EventDelivery, core.DataArray) error); ok {
		r0 = rf(connID, sub, event, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields: ctx, _a1, callbacks
func (_m *Plugin) Init(ctx context.Context, _a1 config.Section, callbacks events.Callbacks) error {
	ret := _m.Called(ctx, _a1, callbacks)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.Section, events.Callbacks) error); ok {
		r0 = rf(ctx, _a1, callbacks)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitConfig provides a mock function with given fields: _a0
func (_m *Plugin) InitConfig(_a0 config.Section) {
	_m.Called(_a0)
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

// ValidateOptions provides a mock function with given fields: options
func (_m *Plugin) ValidateOptions(options *core.SubscriptionOptions) error {
	ret := _m.Called(options)

	var r0 error
	if rf, ok := ret.Get(0).(func(*core.SubscriptionOptions) error); ok {
		r0 = rf(options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
