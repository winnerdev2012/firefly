// Code generated by mockery v1.0.0. DO NOT EDIT.

package definitionsmocks

import (
	context "context"

	core "github.com/hyperledger/firefly/pkg/core"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Sender is an autogenerated mock type for the Sender type
type Sender struct {
	mock.Mock
}

// ClaimIdentity provides a mock function with given fields: ctx, def, signingIdentity, parentSigner, waitConfirm
func (_m *Sender) ClaimIdentity(ctx context.Context, def *core.IdentityClaim, signingIdentity *core.SignerRef, parentSigner *core.SignerRef, waitConfirm bool) error {
	ret := _m.Called(ctx, def, signingIdentity, parentSigner, waitConfirm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.IdentityClaim, *core.SignerRef, *core.SignerRef, bool) error); ok {
		r0 = rf(ctx, def, signingIdentity, parentSigner, waitConfirm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DefineContractAPI provides a mock function with given fields: ctx, httpServerURL, api, waitConfirm
func (_m *Sender) DefineContractAPI(ctx context.Context, httpServerURL string, api *core.ContractAPI, waitConfirm bool) error {
	ret := _m.Called(ctx, httpServerURL, api, waitConfirm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *core.ContractAPI, bool) error); ok {
		r0 = rf(ctx, httpServerURL, api, waitConfirm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DefineDatatype provides a mock function with given fields: ctx, datatype, waitConfirm
func (_m *Sender) DefineDatatype(ctx context.Context, datatype *core.Datatype, waitConfirm bool) error {
	ret := _m.Called(ctx, datatype, waitConfirm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.Datatype, bool) error); ok {
		r0 = rf(ctx, datatype, waitConfirm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DefineFFI provides a mock function with given fields: ctx, ffi, waitConfirm
func (_m *Sender) DefineFFI(ctx context.Context, ffi *fftypes.FFI, waitConfirm bool) error {
	ret := _m.Called(ctx, ffi, waitConfirm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.FFI, bool) error); ok {
		r0 = rf(ctx, ffi, waitConfirm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DefineTokenPool provides a mock function with given fields: ctx, pool, waitConfirm
func (_m *Sender) DefineTokenPool(ctx context.Context, pool *core.TokenPoolAnnouncement, waitConfirm bool) error {
	ret := _m.Called(ctx, pool, waitConfirm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.TokenPoolAnnouncement, bool) error); ok {
		r0 = rf(ctx, pool, waitConfirm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Sender) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// UpdateIdentity provides a mock function with given fields: ctx, identity, def, signingIdentity, waitConfirm
func (_m *Sender) UpdateIdentity(ctx context.Context, identity *core.Identity, def *core.IdentityUpdate, signingIdentity *core.SignerRef, waitConfirm bool) error {
	ret := _m.Called(ctx, identity, def, signingIdentity, waitConfirm)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.Identity, *core.IdentityUpdate, *core.SignerRef, bool) error); ok {
		r0 = rf(ctx, identity, def, signingIdentity, waitConfirm)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
