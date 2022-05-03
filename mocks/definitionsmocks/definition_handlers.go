// Code generated by mockery v1.0.0. DO NOT EDIT.

package definitionsmocks

import (
	context "context"

	definitions "github.com/hyperledger/firefly/internal/definitions"
	fftypes "github.com/hyperledger/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// DefinitionHandlers is an autogenerated mock type for the DefinitionHandlers type
type DefinitionHandlers struct {
	mock.Mock
}

// HandleDefinitionBroadcast provides a mock function with given fields: ctx, state, msg, data, tx
func (_m *DefinitionHandlers) HandleDefinitionBroadcast(ctx context.Context, state definitions.DefinitionBatchState, msg *fftypes.Message, data fftypes.DataArray, tx *fftypes.UUID) (definitions.HandlerResult, error) {
	ret := _m.Called(ctx, state, msg, data, tx)

	var r0 definitions.HandlerResult
	if rf, ok := ret.Get(0).(func(context.Context, definitions.DefinitionBatchState, *fftypes.Message, fftypes.DataArray, *fftypes.UUID) definitions.HandlerResult); ok {
		r0 = rf(ctx, state, msg, data, tx)
	} else {
		r0 = ret.Get(0).(definitions.HandlerResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, definitions.DefinitionBatchState, *fftypes.Message, fftypes.DataArray, *fftypes.UUID) error); ok {
		r1 = rf(ctx, state, msg, data, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendReply provides a mock function with given fields: ctx, event, reply
func (_m *DefinitionHandlers) SendReply(ctx context.Context, event *fftypes.Event, reply *fftypes.MessageInOut) {
	_m.Called(ctx, event, reply)
}
