// Code generated by mockery v1.0.0. DO NOT EDIT.

package networkmapmocks

import (
	context "context"

	database "github.com/hyperledger/firefly/pkg/database"
	fftypes "github.com/hyperledger/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// GetIdentities provides a mock function with given fields: ctx, filter
func (_m *Manager) GetIdentities(ctx context.Context, filter database.AndFilter) ([]*fftypes.Identity, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Identity); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Identity)
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

// GetIdentityByID provides a mock function with given fields: ctx, id
func (_m *Manager) GetIdentityByID(ctx context.Context, id string) (*fftypes.Identity, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Identity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Identity)
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

// GetNodeByID provides a mock function with given fields: ctx, id
func (_m *Manager) GetNodeByID(ctx context.Context, id string) (*fftypes.Identity, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Identity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Identity)
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
func (_m *Manager) GetNodes(ctx context.Context, filter database.AndFilter) ([]*fftypes.Identity, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Identity); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Identity)
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

// GetOrganizationByID provides a mock function with given fields: ctx, id
func (_m *Manager) GetOrganizationByID(ctx context.Context, id string) (*fftypes.Identity, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Identity); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Identity)
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
func (_m *Manager) GetOrganizations(ctx context.Context, filter database.AndFilter) ([]*fftypes.Identity, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Identity); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Identity)
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

// RegisterNode provides a mock function with given fields: ctx, waitConfirm
func (_m *Manager) RegisterNode(ctx context.Context, waitConfirm bool) (*fftypes.Identity, error) {
	ret := _m.Called(ctx, waitConfirm)

	var r0 *fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, bool) *fftypes.Identity); ok {
		r0 = rf(ctx, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Identity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, bool) error); ok {
		r1 = rf(ctx, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterNodeOrganization provides a mock function with given fields: ctx, waitConfirm
func (_m *Manager) RegisterNodeOrganization(ctx context.Context, waitConfirm bool) (*fftypes.Identity, error) {
	ret := _m.Called(ctx, waitConfirm)

	var r0 *fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, bool) *fftypes.Identity); ok {
		r0 = rf(ctx, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Identity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, bool) error); ok {
		r1 = rf(ctx, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterOrganization provides a mock function with given fields: ctx, org, waitConfirm
func (_m *Manager) RegisterOrganization(ctx context.Context, org *fftypes.IdentityCreateDTO, waitConfirm bool) (*fftypes.Identity, error) {
	ret := _m.Called(ctx, org, waitConfirm)

	var r0 *fftypes.Identity
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.IdentityCreateDTO, bool) *fftypes.Identity); ok {
		r0 = rf(ctx, org, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Identity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.IdentityCreateDTO, bool) error); ok {
		r1 = rf(ctx, org, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
