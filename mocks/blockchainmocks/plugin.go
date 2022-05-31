// Code generated by mockery v1.0.0. DO NOT EDIT.

package blockchainmocks

import (
	config "github.com/hyperledger/firefly-common/pkg/config"
	blockchain "github.com/hyperledger/firefly/pkg/blockchain"

	context "context"

	core "github.com/hyperledger/firefly/pkg/core"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	metrics "github.com/hyperledger/firefly/internal/metrics"

	mock "github.com/stretchr/testify/mock"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

// AddContractListener provides a mock function with given fields: ctx, subscription
func (_m *Plugin) AddContractListener(ctx context.Context, subscription *core.ContractListenerInput) error {
	ret := _m.Called(ctx, subscription)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.ContractListenerInput) error); ok {
		r0 = rf(ctx, subscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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

// ConfigureContract provides a mock function with given fields: ctx, contracts
func (_m *Plugin) ConfigureContract(ctx context.Context, contracts *core.FireFlyContracts) error {
	ret := _m.Called(ctx, contracts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.FireFlyContracts) error); ok {
		r0 = rf(ctx, contracts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteContractListener provides a mock function with given fields: ctx, subscription
func (_m *Plugin) DeleteContractListener(ctx context.Context, subscription *core.ContractListener) error {
	ret := _m.Called(ctx, subscription)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.ContractListener) error); ok {
		r0 = rf(ctx, subscription)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateEventSignature provides a mock function with given fields: ctx, event
func (_m *Plugin) GenerateEventSignature(ctx context.Context, event *core.FFIEventDefinition) string {
	ret := _m.Called(ctx, event)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *core.FFIEventDefinition) string); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GenerateFFI provides a mock function with given fields: ctx, generationRequest
func (_m *Plugin) GenerateFFI(ctx context.Context, generationRequest *core.FFIGenerationRequest) (*core.FFI, error) {
	ret := _m.Called(ctx, generationRequest)

	var r0 *core.FFI
	if rf, ok := ret.Get(0).(func(context.Context, *core.FFIGenerationRequest) *core.FFI); ok {
		r0 = rf(ctx, generationRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *core.FFIGenerationRequest) error); ok {
		r1 = rf(ctx, generationRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFFIParamValidator provides a mock function with given fields: ctx
func (_m *Plugin) GetFFIParamValidator(ctx context.Context) (core.FFIParamValidator, error) {
	ret := _m.Called(ctx)

	var r0 core.FFIParamValidator
	if rf, ok := ret.Get(0).(func(context.Context) core.FFIParamValidator); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.FFIParamValidator)
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

// Init provides a mock function with given fields: ctx, _a1, callbacks, _a3
func (_m *Plugin) Init(ctx context.Context, _a1 config.Section, callbacks blockchain.Callbacks, _a3 metrics.Manager) error {
	ret := _m.Called(ctx, _a1, callbacks, _a3)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.Section, blockchain.Callbacks, metrics.Manager) error); ok {
		r0 = rf(ctx, _a1, callbacks, _a3)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitConfig provides a mock function with given fields: _a0
func (_m *Plugin) InitConfig(_a0 config.Section) {
	_m.Called(_a0)
}

// InvokeContract provides a mock function with given fields: ctx, operationID, signingKey, location, method, input, options
func (_m *Plugin) InvokeContract(ctx context.Context, operationID *fftypes.UUID, signingKey string, location *fftypes.JSONAny, method *core.FFIMethod, input map[string]interface{}, options map[string]interface{}) error {
	ret := _m.Called(ctx, operationID, signingKey, location, method, input, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, string, *fftypes.JSONAny, *core.FFIMethod, map[string]interface{}, map[string]interface{}) error); ok {
		r0 = rf(ctx, operationID, signingKey, location, method, input, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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

// NetworkVersion provides a mock function with given fields: ctx
func (_m *Plugin) NetworkVersion(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NormalizeContractLocation provides a mock function with given fields: ctx, location
func (_m *Plugin) NormalizeContractLocation(ctx context.Context, location *fftypes.JSONAny) (*fftypes.JSONAny, error) {
	ret := _m.Called(ctx, location)

	var r0 *fftypes.JSONAny
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.JSONAny) *fftypes.JSONAny); ok {
		r0 = rf(ctx, location)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.JSONAny)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.JSONAny) error); ok {
		r1 = rf(ctx, location)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NormalizeSigningKey provides a mock function with given fields: ctx, keyRef
func (_m *Plugin) NormalizeSigningKey(ctx context.Context, keyRef string) (string, error) {
	ret := _m.Called(ctx, keyRef)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, keyRef)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, keyRef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryContract provides a mock function with given fields: ctx, location, method, input, options
func (_m *Plugin) QueryContract(ctx context.Context, location *fftypes.JSONAny, method *core.FFIMethod, input map[string]interface{}, options map[string]interface{}) (interface{}, error) {
	ret := _m.Called(ctx, location, method, input, options)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.JSONAny, *core.FFIMethod, map[string]interface{}, map[string]interface{}) interface{}); ok {
		r0 = rf(ctx, location, method, input, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.JSONAny, *core.FFIMethod, map[string]interface{}, map[string]interface{}) error); ok {
		r1 = rf(ctx, location, method, input, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

// SubmitBatchPin provides a mock function with given fields: ctx, operationID, signingKey, batch
func (_m *Plugin) SubmitBatchPin(ctx context.Context, operationID *fftypes.UUID, signingKey string, batch *blockchain.BatchPin) error {
	ret := _m.Called(ctx, operationID, signingKey, batch)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, string, *blockchain.BatchPin) error); ok {
		r0 = rf(ctx, operationID, signingKey, batch)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubmitNetworkAction provides a mock function with given fields: ctx, operationID, signingKey, action
func (_m *Plugin) SubmitNetworkAction(ctx context.Context, operationID *fftypes.UUID, signingKey string, action core.FFEnum) error {
	ret := _m.Called(ctx, operationID, signingKey, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, string, core.FFEnum) error); ok {
		r0 = rf(ctx, operationID, signingKey, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TerminateContract provides a mock function with given fields: ctx, contracts, termination
func (_m *Plugin) TerminateContract(ctx context.Context, contracts *core.FireFlyContracts, termination *blockchain.Event) error {
	ret := _m.Called(ctx, contracts, termination)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *core.FireFlyContracts, *blockchain.Event) error); ok {
		r0 = rf(ctx, contracts, termination)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifierType provides a mock function with given fields:
func (_m *Plugin) VerifierType() core.FFEnum {
	ret := _m.Called()

	var r0 core.FFEnum
	if rf, ok := ret.Get(0).(func() core.FFEnum); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(core.FFEnum)
	}

	return r0
}
