// Code generated by mockery v1.0.0. DO NOT EDIT.

package syshandlersmocks

import (
	context "context"

	database "github.com/hyperledger-labs/firefly/pkg/database"
	fftypes "github.com/hyperledger-labs/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// SystemHandlers is an autogenerated mock type for the SystemHandlers type
type SystemHandlers struct {
	mock.Mock
}

// EnsureLocalGroup provides a mock function with given fields: ctx, group
func (_m *SystemHandlers) EnsureLocalGroup(ctx context.Context, group *fftypes.Group) (bool, error) {
	ret := _m.Called(ctx, group)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.Group) bool); ok {
		r0 = rf(ctx, group)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.Group) error); ok {
		r1 = rf(ctx, group)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGroupByID provides a mock function with given fields: ctx, id
func (_m *SystemHandlers) GetGroupByID(ctx context.Context, id string) (*fftypes.Group, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.Group
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Group); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGroups provides a mock function with given fields: ctx, filter
func (_m *SystemHandlers) GetGroups(ctx context.Context, filter database.AndFilter) ([]*fftypes.Group, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Group
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Group); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Group)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, database.AndFilter) error); ok {
		r2 = rf(ctx, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HandleSystemBroadcast provides a mock function with given fields: ctx, msg, data
func (_m *SystemHandlers) HandleSystemBroadcast(ctx context.Context, msg *fftypes.Message, data []*fftypes.Data) (bool, error) {
	ret := _m.Called(ctx, msg, data)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.Message, []*fftypes.Data) bool); ok {
		r0 = rf(ctx, msg, data)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.Message, []*fftypes.Data) error); ok {
		r1 = rf(ctx, msg, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResolveInitGroup provides a mock function with given fields: ctx, msg
func (_m *SystemHandlers) ResolveInitGroup(ctx context.Context, msg *fftypes.Message) (*fftypes.Group, error) {
	ret := _m.Called(ctx, msg)

	var r0 *fftypes.Group
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.Message) *fftypes.Group); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Group)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.Message) error); ok {
		r1 = rf(ctx, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendReply provides a mock function with given fields: ctx, event, reply
func (_m *SystemHandlers) SendReply(ctx context.Context, event *fftypes.Event, reply *fftypes.MessageInOut) {
	_m.Called(ctx, event, reply)
}
