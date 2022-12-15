// Code generated by mockery v2.15.0. DO NOT EDIT.

package tokenmocks

import (
	tokens "github.com/hyperledger/firefly/pkg/tokens"
	mock "github.com/stretchr/testify/mock"
)

// Callbacks is an autogenerated mock type for the Callbacks type
type Callbacks struct {
	mock.Mock
}

// TokenPoolCreated provides a mock function with given fields: plugin, pool
func (_m *Callbacks) TokenPoolCreated(plugin tokens.Plugin, pool *tokens.TokenPool) error {
	ret := _m.Called(plugin, pool)

	var r0 error
	if rf, ok := ret.Get(0).(func(tokens.Plugin, *tokens.TokenPool) error); ok {
		r0 = rf(plugin, pool)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TokensApproved provides a mock function with given fields: plugin, approval
func (_m *Callbacks) TokensApproved(plugin tokens.Plugin, approval *tokens.TokenApproval) error {
	ret := _m.Called(plugin, approval)

	var r0 error
	if rf, ok := ret.Get(0).(func(tokens.Plugin, *tokens.TokenApproval) error); ok {
		r0 = rf(plugin, approval)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TokensTransferred provides a mock function with given fields: plugin, transfer
func (_m *Callbacks) TokensTransferred(plugin tokens.Plugin, transfer *tokens.TokenTransfer) error {
	ret := _m.Called(plugin, transfer)

	var r0 error
	if rf, ok := ret.Get(0).(func(tokens.Plugin, *tokens.TokenTransfer) error); ok {
		r0 = rf(plugin, transfer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCallbacks interface {
	mock.TestingT
	Cleanup(func())
}

// NewCallbacks creates a new instance of Callbacks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCallbacks(t mockConstructorTestingTNewCallbacks) *Callbacks {
	mock := &Callbacks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
