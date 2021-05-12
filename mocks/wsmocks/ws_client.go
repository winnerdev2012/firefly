// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package wsmocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// WSClient is an autogenerated mock type for the WSClient type
type WSClient struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *WSClient) Close() {
	_m.Called()
}

// Receive provides a mock function with given fields:
func (_m *WSClient) Receive() <-chan []byte {
	ret := _m.Called()

	var r0 <-chan []byte
	if rf, ok := ret.Get(0).(func() <-chan []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan []byte)
		}
	}

	return r0
}

// Send provides a mock function with given fields: ctx, message
func (_m *WSClient) Send(ctx context.Context, message []byte) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
