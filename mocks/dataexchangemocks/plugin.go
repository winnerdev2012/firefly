// Code generated by mockery v1.0.0. DO NOT EDIT.

package dataexchangemocks

import (
	context "context"

	config "github.com/hyperledger/firefly-common/pkg/config"

	dataexchange "github.com/hyperledger/firefly/pkg/dataexchange"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

// AddPeer provides a mock function with given fields: ctx, peer
func (_m *Plugin) AddPeer(ctx context.Context, peer fftypes.JSONObject) error {
	ret := _m.Called(ctx, peer)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, fftypes.JSONObject) error); ok {
		r0 = rf(ctx, peer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Capabilities provides a mock function with given fields:
func (_m *Plugin) Capabilities() *dataexchange.Capabilities {
	ret := _m.Called()

	var r0 *dataexchange.Capabilities
	if rf, ok := ret.Get(0).(func() *dataexchange.Capabilities); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dataexchange.Capabilities)
		}
	}

	return r0
}

// CheckBlobReceived provides a mock function with given fields: ctx, peerID, ns, id
func (_m *Plugin) CheckBlobReceived(ctx context.Context, peerID string, ns string, id fftypes.UUID) (*fftypes.Bytes32, int64, error) {
	ret := _m.Called(ctx, peerID, ns, id)

	var r0 *fftypes.Bytes32
	if rf, ok := ret.Get(0).(func(context.Context, string, string, fftypes.UUID) *fftypes.Bytes32); ok {
		r0 = rf(ctx, peerID, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Bytes32)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, string, fftypes.UUID) int64); ok {
		r1 = rf(ctx, peerID, ns, id)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, fftypes.UUID) error); ok {
		r2 = rf(ctx, peerID, ns, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DownloadBlob provides a mock function with given fields: ctx, payloadRef
func (_m *Plugin) DownloadBlob(ctx context.Context, payloadRef string) (io.ReadCloser, error) {
	ret := _m.Called(ctx, payloadRef)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(context.Context, string) io.ReadCloser); ok {
		r0 = rf(ctx, payloadRef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, payloadRef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEndpointInfo provides a mock function with given fields: ctx
func (_m *Plugin) GetEndpointInfo(ctx context.Context) (fftypes.JSONObject, error) {
	ret := _m.Called(ctx)

	var r0 fftypes.JSONObject
	if rf, ok := ret.Get(0).(func(context.Context) fftypes.JSONObject); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fftypes.JSONObject)
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

// Init provides a mock function with given fields: ctx, _a1
func (_m *Plugin) Init(ctx context.Context, _a1 config.Section) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, config.Section) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitConfig provides a mock function with given fields: _a0
func (_m *Plugin) InitConfig(_a0 config.Section) {
	_m.Called(_a0)
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

// SendMessage provides a mock function with given fields: ctx, nsOpID, peerID, data
func (_m *Plugin) SendMessage(ctx context.Context, nsOpID string, peerID string, data []byte) error {
	ret := _m.Called(ctx, nsOpID, peerID, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []byte) error); ok {
		r0 = rf(ctx, nsOpID, peerID, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetHandler provides a mock function with given fields: namespace, handler
func (_m *Plugin) SetHandler(namespace string, handler dataexchange.Callbacks) {
	_m.Called(namespace, handler)
}

// SetNodes provides a mock function with given fields: nodes
func (_m *Plugin) SetNodes(nodes []fftypes.JSONObject) {
	_m.Called(nodes)
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

// TransferBlob provides a mock function with given fields: ctx, nsOpID, peerID, payloadRef
func (_m *Plugin) TransferBlob(ctx context.Context, nsOpID string, peerID string, payloadRef string) error {
	ret := _m.Called(ctx, nsOpID, peerID, payloadRef)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, nsOpID, peerID, payloadRef)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UploadBlob provides a mock function with given fields: ctx, ns, id, content
func (_m *Plugin) UploadBlob(ctx context.Context, ns string, id fftypes.UUID, content io.Reader) (string, *fftypes.Bytes32, int64, error) {
	ret := _m.Called(ctx, ns, id, content)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, fftypes.UUID, io.Reader) string); ok {
		r0 = rf(ctx, ns, id, content)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 *fftypes.Bytes32
	if rf, ok := ret.Get(1).(func(context.Context, string, fftypes.UUID, io.Reader) *fftypes.Bytes32); ok {
		r1 = rf(ctx, ns, id, content)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*fftypes.Bytes32)
		}
	}

	var r2 int64
	if rf, ok := ret.Get(2).(func(context.Context, string, fftypes.UUID, io.Reader) int64); ok {
		r2 = rf(ctx, ns, id, content)
	} else {
		r2 = ret.Get(2).(int64)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, string, fftypes.UUID, io.Reader) error); ok {
		r3 = rf(ctx, ns, id, content)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}
