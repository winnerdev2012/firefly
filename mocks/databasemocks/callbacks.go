// Code generated by mockery v2.16.0. DO NOT EDIT.

package databasemocks

import (
	core "github.com/hyperledger/firefly/pkg/core"
	database "github.com/hyperledger/firefly/pkg/database"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// Callbacks is an autogenerated mock type for the Callbacks type
type Callbacks struct {
	mock.Mock
}

// HashCollectionNSEvent provides a mock function with given fields: resType, eventType, namespace, hash
func (_m *Callbacks) HashCollectionNSEvent(resType database.HashCollectionNS, eventType core.ChangeEventType, namespace string, hash *fftypes.Bytes32) {
	_m.Called(resType, eventType, namespace, hash)
}

// OrderedCollectionNSEvent provides a mock function with given fields: resType, eventType, namespace, sequence
func (_m *Callbacks) OrderedCollectionNSEvent(resType database.OrderedCollectionNS, eventType core.ChangeEventType, namespace string, sequence int64) {
	_m.Called(resType, eventType, namespace, sequence)
}

// OrderedUUIDCollectionNSEvent provides a mock function with given fields: resType, eventType, namespace, id, sequence
func (_m *Callbacks) OrderedUUIDCollectionNSEvent(resType database.OrderedUUIDCollectionNS, eventType core.ChangeEventType, namespace string, id *fftypes.UUID, sequence int64) {
	_m.Called(resType, eventType, namespace, id, sequence)
}

// UUIDCollectionNSEvent provides a mock function with given fields: resType, eventType, namespace, id
func (_m *Callbacks) UUIDCollectionNSEvent(resType database.UUIDCollectionNS, eventType core.ChangeEventType, namespace string, id *fftypes.UUID) {
	_m.Called(resType, eventType, namespace, id)
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
