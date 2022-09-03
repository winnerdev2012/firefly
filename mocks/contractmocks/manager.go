// Code generated by mockery v1.0.0. DO NOT EDIT.

package contractmocks

import (
	context "context"

	core "github.com/hyperledger/firefly/pkg/core"

	database "github.com/hyperledger/firefly/pkg/database"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// AddContractAPIListener provides a mock function with given fields: ctx, apiName, eventPath, listener
func (_m *Manager) AddContractAPIListener(ctx context.Context, apiName string, eventPath string, listener *core.ContractListener) (*core.ContractListener, error) {
	ret := _m.Called(ctx, apiName, eventPath, listener)

	var r0 *core.ContractListener
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *core.ContractListener) *core.ContractListener); ok {
		r0 = rf(ctx, apiName, eventPath, listener)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.ContractListener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *core.ContractListener) error); ok {
		r1 = rf(ctx, apiName, eventPath, listener)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddContractListener provides a mock function with given fields: ctx, listener
func (_m *Manager) AddContractListener(ctx context.Context, listener *core.ContractListenerInput) (*core.ContractListener, error) {
	ret := _m.Called(ctx, listener)

	var r0 *core.ContractListener
	if rf, ok := ret.Get(0).(func(context.Context, *core.ContractListenerInput) *core.ContractListener); ok {
		r0 = rf(ctx, listener)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.ContractListener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *core.ContractListenerInput) error); ok {
		r1 = rf(ctx, listener)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteContractListenerByNameOrID provides a mock function with given fields: ctx, nameOrID
func (_m *Manager) DeleteContractListenerByNameOrID(ctx context.Context, nameOrID string) error {
	ret := _m.Called(ctx, nameOrID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, nameOrID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateFFI provides a mock function with given fields: ctx, generationRequest
func (_m *Manager) GenerateFFI(ctx context.Context, generationRequest *fftypes.FFIGenerationRequest) (*fftypes.FFI, error) {
	ret := _m.Called(ctx, generationRequest)

	var r0 *fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.FFIGenerationRequest) *fftypes.FFI); ok {
		r0 = rf(ctx, generationRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.FFIGenerationRequest) error); ok {
		r1 = rf(ctx, generationRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContractAPI provides a mock function with given fields: ctx, httpServerURL, apiName
func (_m *Manager) GetContractAPI(ctx context.Context, httpServerURL string, apiName string) (*core.ContractAPI, error) {
	ret := _m.Called(ctx, httpServerURL, apiName)

	var r0 *core.ContractAPI
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.ContractAPI); ok {
		r0 = rf(ctx, httpServerURL, apiName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.ContractAPI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, httpServerURL, apiName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContractAPIInterface provides a mock function with given fields: ctx, apiName
func (_m *Manager) GetContractAPIInterface(ctx context.Context, apiName string) (*fftypes.FFI, error) {
	ret := _m.Called(ctx, apiName)

	var r0 *fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.FFI); ok {
		r0 = rf(ctx, apiName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, apiName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContractAPIListeners provides a mock function with given fields: ctx, apiName, eventPath, filter
func (_m *Manager) GetContractAPIListeners(ctx context.Context, apiName string, eventPath string, filter database.AndFilter) ([]*core.ContractListener, *database.FilterResult, error) {
	ret := _m.Called(ctx, apiName, eventPath, filter)

	var r0 []*core.ContractListener
	if rf, ok := ret.Get(0).(func(context.Context, string, string, database.AndFilter) []*core.ContractListener); ok {
		r0 = rf(ctx, apiName, eventPath, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.ContractListener)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, apiName, eventPath, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, database.AndFilter) error); ok {
		r2 = rf(ctx, apiName, eventPath, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetContractAPIs provides a mock function with given fields: ctx, httpServerURL, filter
func (_m *Manager) GetContractAPIs(ctx context.Context, httpServerURL string, filter database.AndFilter) ([]*core.ContractAPI, *database.FilterResult, error) {
	ret := _m.Called(ctx, httpServerURL, filter)

	var r0 []*core.ContractAPI
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.ContractAPI); ok {
		r0 = rf(ctx, httpServerURL, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.ContractAPI)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, httpServerURL, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, database.AndFilter) error); ok {
		r2 = rf(ctx, httpServerURL, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetContractListenerByNameOrID provides a mock function with given fields: ctx, nameOrID
func (_m *Manager) GetContractListenerByNameOrID(ctx context.Context, nameOrID string) (*core.ContractListener, error) {
	ret := _m.Called(ctx, nameOrID)

	var r0 *core.ContractListener
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.ContractListener); ok {
		r0 = rf(ctx, nameOrID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.ContractListener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, nameOrID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContractListenerByNameOrIDWithStatus provides a mock function with given fields: ctx, nameOrID
func (_m *Manager) GetContractListenerByNameOrIDWithStatus(ctx context.Context, nameOrID string) (*core.ContractListenerWithStatus, error) {
	ret := _m.Called(ctx, nameOrID)

	var r0 *core.ContractListenerWithStatus
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.ContractListenerWithStatus); ok {
		r0 = rf(ctx, nameOrID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.ContractListenerWithStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, nameOrID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContractListeners provides a mock function with given fields: ctx, filter
func (_m *Manager) GetContractListeners(ctx context.Context, filter database.AndFilter) ([]*core.ContractListener, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*core.ContractListener
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*core.ContractListener); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.ContractListener)
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

// GetFFI provides a mock function with given fields: ctx, name, version
func (_m *Manager) GetFFI(ctx context.Context, name string, version string) (*fftypes.FFI, error) {
	ret := _m.Called(ctx, name, version)

	var r0 *fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.FFI); ok {
		r0 = rf(ctx, name, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFFIByID provides a mock function with given fields: ctx, id
func (_m *Manager) GetFFIByID(ctx context.Context, id *fftypes.UUID) (*fftypes.FFI, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *fftypes.FFI); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFFIByIDWithChildren provides a mock function with given fields: ctx, id
func (_m *Manager) GetFFIByIDWithChildren(ctx context.Context, id *fftypes.UUID) (*fftypes.FFI, error) {
	ret := _m.Called(ctx, id)

	var r0 *fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *fftypes.FFI); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFFIWithChildren provides a mock function with given fields: ctx, name, version
func (_m *Manager) GetFFIWithChildren(ctx context.Context, name string, version string) (*fftypes.FFI, error) {
	ret := _m.Called(ctx, name, version)

	var r0 *fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.FFI); ok {
		r0 = rf(ctx, name, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.FFI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, name, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFFIs provides a mock function with given fields: ctx, filter
func (_m *Manager) GetFFIs(ctx context.Context, filter database.AndFilter) ([]*fftypes.FFI, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.FFI
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.FFI); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.FFI)
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

// InvokeContract provides a mock function with given fields: ctx, req, waitConfirm
func (_m *Manager) InvokeContract(ctx context.Context, req *core.ContractCallRequest, waitConfirm bool) (interface{}, error) {
	ret := _m.Called(ctx, req, waitConfirm)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, *core.ContractCallRequest, bool) interface{}); ok {
		r0 = rf(ctx, req, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *core.ContractCallRequest, bool) error); ok {
		r1 = rf(ctx, req, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvokeContractAPI provides a mock function with given fields: ctx, apiName, methodPath, req, waitConfirm
func (_m *Manager) InvokeContractAPI(ctx context.Context, apiName string, methodPath string, req *core.ContractCallRequest, waitConfirm bool) (interface{}, error) {
	ret := _m.Called(ctx, apiName, methodPath, req, waitConfirm)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *core.ContractCallRequest, bool) interface{}); ok {
		r0 = rf(ctx, apiName, methodPath, req, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, *core.ContractCallRequest, bool) error); ok {
		r1 = rf(ctx, apiName, methodPath, req, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields:
func (_m *Manager) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
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

// ResolveContractAPI provides a mock function with given fields: ctx, httpServerURL, api
func (_m *Manager) ResolveContractAPI(ctx context.Context, httpServerURL string, api *core.ContractAPI) error {
	ret := _m.Called(ctx, httpServerURL, api)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *core.ContractAPI) error); ok {
		r0 = rf(ctx, httpServerURL, api)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResolveFFI provides a mock function with given fields: ctx, ffi
func (_m *Manager) ResolveFFI(ctx context.Context, ffi *fftypes.FFI) error {
	ret := _m.Called(ctx, ffi)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.FFI) error); ok {
		r0 = rf(ctx, ffi)
	} else {
		r0 = ret.Error(0)
	}

	return r0
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
