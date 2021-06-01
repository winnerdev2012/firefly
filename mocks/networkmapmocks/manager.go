// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package networkmapmocks

import (
	context "context"

	database "github.com/kaleido-io/firefly/pkg/database"
	fftypes "github.com/kaleido-io/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// GetNodeByID provides a mock function with given fields: ctx, id
func (_m *Manager) GetNodeByID(ctx context.Context, id string) (*fftypes.Node, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.Node
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Node); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Node)
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

// GetNodes provides a mock function with given fields: ctx, filter
func (_m *Manager) GetNodes(ctx context.Context, filter database.AndFilter) ([]*fftypes.Node, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Node
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Node); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, database.AndFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrganizationByID provides a mock function with given fields: ctx, id
func (_m *Manager) GetOrganizationByID(ctx context.Context, id string) (*fftypes.Organization, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.Organization
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Organization); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Organization)
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

// GetOrganizations provides a mock function with given fields: ctx, filter
func (_m *Manager) GetOrganizations(ctx context.Context, filter database.AndFilter) ([]*fftypes.Organization, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Organization
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Organization); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Organization)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, database.AndFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterNode provides a mock function with given fields: ctx
func (_m *Manager) RegisterNode(ctx context.Context) (*fftypes.Message, error) {
	ret := _m.Called(ctx)

	var r0 *fftypes.Message
	if rf, ok := ret.Get(0).(func(context.Context) *fftypes.Message); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterOrganization provides a mock function with given fields: ctx, org
func (_m *Manager) RegisterOrganization(ctx context.Context, org *fftypes.Organization) (*fftypes.Message, error) {
	ret := _m.Called(ctx, org)

	var r0 *fftypes.Message
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.Organization) *fftypes.Message); ok {
		r0 = rf(ctx, org)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.Organization) error); ok {
		r1 = rf(ctx, org)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
