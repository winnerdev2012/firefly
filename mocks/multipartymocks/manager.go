// Code generated by mockery v1.0.0. DO NOT EDIT.

package multipartymocks

import (
	context "context"

	blockchain "github.com/hyperledger/firefly/pkg/blockchain"

	core "github.com/hyperledger/firefly/pkg/core"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// ConfigureContract provides a mock function with given fields: ctx, contracts
func (_m *Manager) ConfigureContract(ctx context.Context, contracts *core.FireFlyContracts) error {
	ret := _m.Called(ctx, contracts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.FireFlyContracts) error); ok {
		r0 = rf(ctx, contracts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetNetworkVersion provides a mock function with given fields:
func (_m *Manager) GetNetworkVersion() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// PrepareOperation provides a mock function with given fields: ctx, op
func (_m *Manager) PrepareOperation(ctx context.Context, op *core.Operation) (*core.PreparedOperation, error) {
	ret := _m.Called(ctx, op)

	var r0 *core.PreparedOperation
	if rf, ok := ret.Get(0).(func(context.Context, *core.Operation) *core.PreparedOperation); ok {
		r0 = rf(ctx, op)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.PreparedOperation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *core.Operation) error); ok {
		r1 = rf(ctx, op)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunOperation provides a mock function with given fields: ctx, op
func (_m *Manager) RunOperation(ctx context.Context, op *core.PreparedOperation) (fftypes.JSONObject, bool, error) {
	ret := _m.Called(ctx, op)

	var r0 fftypes.JSONObject
	if rf, ok := ret.Get(0).(func(context.Context, *core.PreparedOperation) fftypes.JSONObject); ok {
		r0 = rf(ctx, op)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fftypes.JSONObject)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, *core.PreparedOperation) bool); ok {
		r1 = rf(ctx, op)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *core.PreparedOperation) error); ok {
		r2 = rf(ctx, op)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SubmitBatchPin provides a mock function with given fields: ctx, batch, contexts, payloadRef
func (_m *Manager) SubmitBatchPin(ctx context.Context, batch *core.BatchPersisted, contexts []*fftypes.Bytes32, payloadRef string) error {
	ret := _m.Called(ctx, batch, contexts, payloadRef)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.BatchPersisted, []*fftypes.Bytes32, string) error); ok {
		r0 = rf(ctx, batch, contexts, payloadRef)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitNetworkAction provides a mock function with given fields: ctx, nsOpID, signingKey, action
func (_m *Manager) SubmitNetworkAction(ctx context.Context, nsOpID string, signingKey string, action fftypes.FFEnum) error {
	ret := _m.Called(ctx, nsOpID, signingKey, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, fftypes.FFEnum) error); ok {
		r0 = rf(ctx, nsOpID, signingKey, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TerminateContract provides a mock function with given fields: ctx, contracts, termination
func (_m *Manager) TerminateContract(ctx context.Context, contracts *core.FireFlyContracts, termination *blockchain.Event) error {
	ret := _m.Called(ctx, contracts, termination)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.FireFlyContracts, *blockchain.Event) error); ok {
		r0 = rf(ctx, contracts, termination)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
