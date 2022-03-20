// Code generated by mockery v1.0.0. DO NOT EDIT.

package assetmocks

import (
	context "context"

	database "github.com/hyperledger/firefly/pkg/database"
	fftypes "github.com/hyperledger/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"

	sysmessaging "github.com/hyperledger/firefly/internal/sysmessaging"
)

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// ActivateTokenPool provides a mock function with given fields: ctx, pool, blockchainInfo
func (_m *Manager) ActivateTokenPool(ctx context.Context, pool *fftypes.TokenPool, blockchainInfo fftypes.JSONObject) error {
	ret := _m.Called(ctx, pool, blockchainInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.TokenPool, fftypes.JSONObject) error); ok {
		r0 = rf(ctx, pool, blockchainInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// BurnTokens provides a mock function with given fields: ctx, ns, transfer, waitConfirm
func (_m *Manager) BurnTokens(ctx context.Context, ns string, transfer *fftypes.TokenTransferInput, waitConfirm bool) (*fftypes.TokenTransfer, error) {
	ret := _m.Called(ctx, ns, transfer, waitConfirm)

	var r0 *fftypes.TokenTransfer
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.TokenTransferInput, bool) *fftypes.TokenTransfer); ok {
		r0 = rf(ctx, ns, transfer, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenTransfer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.TokenTransferInput, bool) error); ok {
		r1 = rf(ctx, ns, transfer, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTokenPool provides a mock function with given fields: ctx, ns, pool, waitConfirm
func (_m *Manager) CreateTokenPool(ctx context.Context, ns string, pool *fftypes.TokenPool, waitConfirm bool) (*fftypes.TokenPool, error) {
	ret := _m.Called(ctx, ns, pool, waitConfirm)

	var r0 *fftypes.TokenPool
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.TokenPool, bool) *fftypes.TokenPool); ok {
		r0 = rf(ctx, ns, pool, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenPool)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.TokenPool, bool) error); ok {
		r1 = rf(ctx, ns, pool, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenAccountPools provides a mock function with given fields: ctx, ns, key, filter
func (_m *Manager) GetTokenAccountPools(ctx context.Context, ns string, key string, filter database.AndFilter) ([]*fftypes.TokenAccountPool, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, key, filter)

	var r0 []*fftypes.TokenAccountPool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, database.AndFilter) []*fftypes.TokenAccountPool); ok {
		r0 = rf(ctx, ns, key, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenAccountPool)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, key, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, key, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTokenAccounts provides a mock function with given fields: ctx, ns, filter
func (_m *Manager) GetTokenAccounts(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.TokenAccount, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.TokenAccount
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.TokenAccount); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenAccount)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTokenApprovals provides a mock function with given fields: ctx, ns, filter
func (_m *Manager) GetTokenApprovals(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.TokenApproval, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.TokenApproval
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.TokenApproval); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenApproval)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTokenBalances provides a mock function with given fields: ctx, ns, filter
func (_m *Manager) GetTokenBalances(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.TokenBalance, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.TokenBalance
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.TokenBalance); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenBalance)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTokenConnectors provides a mock function with given fields: ctx, ns
func (_m *Manager) GetTokenConnectors(ctx context.Context, ns string) ([]*fftypes.TokenConnector, error) {
	ret := _m.Called(ctx, ns)

	var r0 []*fftypes.TokenConnector
	if rf, ok := ret.Get(0).(func(context.Context, string) []*fftypes.TokenConnector); ok {
		r0 = rf(ctx, ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenConnector)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ns)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenPool provides a mock function with given fields: ctx, ns, connector, poolName
func (_m *Manager) GetTokenPool(ctx context.Context, ns string, connector string, poolName string) (*fftypes.TokenPool, error) {
	ret := _m.Called(ctx, ns, connector, poolName)

	var r0 *fftypes.TokenPool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *fftypes.TokenPool); ok {
		r0 = rf(ctx, ns, connector, poolName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenPool)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, ns, connector, poolName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenPoolByNameOrID provides a mock function with given fields: ctx, ns, poolNameOrID
func (_m *Manager) GetTokenPoolByNameOrID(ctx context.Context, ns string, poolNameOrID string) (*fftypes.TokenPool, error) {
	ret := _m.Called(ctx, ns, poolNameOrID)

	var r0 *fftypes.TokenPool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.TokenPool); ok {
		r0 = rf(ctx, ns, poolNameOrID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenPool)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, ns, poolNameOrID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenPools provides a mock function with given fields: ctx, ns, filter
func (_m *Manager) GetTokenPools(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.TokenPool, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.TokenPool
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.TokenPool); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenPool)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTokenTransferByID provides a mock function with given fields: ctx, ns, id
func (_m *Manager) GetTokenTransferByID(ctx context.Context, ns string, id string) (*fftypes.TokenTransfer, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.TokenTransfer
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.TokenTransfer); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenTransfer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, ns, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTokenTransfers provides a mock function with given fields: ctx, ns, filter
func (_m *Manager) GetTokenTransfers(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.TokenTransfer, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.TokenTransfer
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.TokenTransfer); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.TokenTransfer)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MintTokens provides a mock function with given fields: ctx, ns, transfer, waitConfirm
func (_m *Manager) MintTokens(ctx context.Context, ns string, transfer *fftypes.TokenTransferInput, waitConfirm bool) (*fftypes.TokenTransfer, error) {
	ret := _m.Called(ctx, ns, transfer, waitConfirm)

	var r0 *fftypes.TokenTransfer
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.TokenTransferInput, bool) *fftypes.TokenTransfer); ok {
		r0 = rf(ctx, ns, transfer, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenTransfer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.TokenTransferInput, bool) error); ok {
		r1 = rf(ctx, ns, transfer, waitConfirm)
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

// NewApproval provides a mock function with given fields: ns, approve
func (_m *Manager) NewApproval(ns string, approve *fftypes.TokenApprovalInput) sysmessaging.MessageSender {
	ret := _m.Called(ns, approve)

	var r0 sysmessaging.MessageSender
	if rf, ok := ret.Get(0).(func(string, *fftypes.TokenApprovalInput) sysmessaging.MessageSender); ok {
		r0 = rf(ns, approve)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sysmessaging.MessageSender)
		}
	}

	return r0
}

// NewTransfer provides a mock function with given fields: ns, transfer
func (_m *Manager) NewTransfer(ns string, transfer *fftypes.TokenTransferInput) sysmessaging.MessageSender {
	ret := _m.Called(ns, transfer)

	var r0 sysmessaging.MessageSender
	if rf, ok := ret.Get(0).(func(string, *fftypes.TokenTransferInput) sysmessaging.MessageSender); ok {
		r0 = rf(ns, transfer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sysmessaging.MessageSender)
		}
	}

	return r0
}

// PrepareOperation provides a mock function with given fields: ctx, op
func (_m *Manager) PrepareOperation(ctx context.Context, op *fftypes.Operation) (*fftypes.PreparedOperation, error) {
	ret := _m.Called(ctx, op)

	var r0 *fftypes.PreparedOperation
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.Operation) *fftypes.PreparedOperation); ok {
		r0 = rf(ctx, op)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.PreparedOperation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.Operation) error); ok {
		r1 = rf(ctx, op)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunOperation provides a mock function with given fields: ctx, op
func (_m *Manager) RunOperation(ctx context.Context, op *fftypes.PreparedOperation) (fftypes.JSONObject, bool, error) {
	ret := _m.Called(ctx, op)

	var r0 fftypes.JSONObject
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.PreparedOperation) fftypes.JSONObject); ok {
		r0 = rf(ctx, op)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fftypes.JSONObject)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.PreparedOperation) bool); ok {
		r1 = rf(ctx, op)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, *fftypes.PreparedOperation) error); ok {
		r2 = rf(ctx, op)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TokenApproval provides a mock function with given fields: ctx, ns, approval, waitConfirm
func (_m *Manager) TokenApproval(ctx context.Context, ns string, approval *fftypes.TokenApprovalInput, waitConfirm bool) (*fftypes.TokenApproval, error) {
	ret := _m.Called(ctx, ns, approval, waitConfirm)

	var r0 *fftypes.TokenApproval
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.TokenApprovalInput, bool) *fftypes.TokenApproval); ok {
		r0 = rf(ctx, ns, approval, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenApproval)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.TokenApprovalInput, bool) error); ok {
		r1 = rf(ctx, ns, approval, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferTokens provides a mock function with given fields: ctx, ns, transfer, waitConfirm
func (_m *Manager) TransferTokens(ctx context.Context, ns string, transfer *fftypes.TokenTransferInput, waitConfirm bool) (*fftypes.TokenTransfer, error) {
	ret := _m.Called(ctx, ns, transfer, waitConfirm)

	var r0 *fftypes.TokenTransfer
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.TokenTransferInput, bool) *fftypes.TokenTransfer); ok {
		r0 = rf(ctx, ns, transfer, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.TokenTransfer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.TokenTransferInput, bool) error); ok {
		r1 = rf(ctx, ns, transfer, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
