// Code generated by mockery v1.0.0. DO NOT EDIT.

package orchestratormocks

import (
	assets "github.com/hyperledger/firefly/internal/assets"
	broadcast "github.com/hyperledger/firefly/internal/broadcast"

	context "context"

	contracts "github.com/hyperledger/firefly/internal/contracts"

	data "github.com/hyperledger/firefly/internal/data"

	database "github.com/hyperledger/firefly/pkg/database"

	events "github.com/hyperledger/firefly/internal/events"

	fftypes "github.com/hyperledger/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"

	networkmap "github.com/hyperledger/firefly/internal/networkmap"

	privatemessaging "github.com/hyperledger/firefly/internal/privatemessaging"
)

// Orchestrator is an autogenerated mock type for the Orchestrator type
type Orchestrator struct {
	mock.Mock
}

// Assets provides a mock function with given fields:
func (_m *Orchestrator) Assets() assets.Manager {
	ret := _m.Called()

	var r0 assets.Manager
	if rf, ok := ret.Get(0).(func() assets.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(assets.Manager)
		}
	}

	return r0
}

// Broadcast provides a mock function with given fields:
func (_m *Orchestrator) Broadcast() broadcast.Manager {
	ret := _m.Called()

	var r0 broadcast.Manager
	if rf, ok := ret.Get(0).(func() broadcast.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(broadcast.Manager)
		}
	}

	return r0
}

// Contracts provides a mock function with given fields:
func (_m *Orchestrator) Contracts() contracts.Manager {
	ret := _m.Called()

	var r0 contracts.Manager
	if rf, ok := ret.Get(0).(func() contracts.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(contracts.Manager)
		}
	}

	return r0
}

// CreateSubscription provides a mock function with given fields: ctx, ns, subDef
func (_m *Orchestrator) CreateSubscription(ctx context.Context, ns string, subDef *fftypes.Subscription) (*fftypes.Subscription, error) {
	ret := _m.Called(ctx, ns, subDef)

	var r0 *fftypes.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.Subscription) *fftypes.Subscription); ok {
		r0 = rf(ctx, ns, subDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.Subscription) error); ok {
		r1 = rf(ctx, ns, subDef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUpdateSubscription provides a mock function with given fields: ctx, ns, subDef
func (_m *Orchestrator) CreateUpdateSubscription(ctx context.Context, ns string, subDef *fftypes.Subscription) (*fftypes.Subscription, error) {
	ret := _m.Called(ctx, ns, subDef)

	var r0 *fftypes.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.Subscription) *fftypes.Subscription); ok {
		r0 = rf(ctx, ns, subDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.Subscription) error); ok {
		r1 = rf(ctx, ns, subDef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Data provides a mock function with given fields:
func (_m *Orchestrator) Data() data.Manager {
	ret := _m.Called()

	var r0 data.Manager
	if rf, ok := ret.Get(0).(func() data.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(data.Manager)
		}
	}

	return r0
}

// DeleteConfigRecord provides a mock function with given fields: ctx, key
func (_m *Orchestrator) DeleteConfigRecord(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSubscription provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) DeleteSubscription(ctx context.Context, ns string, id string) error {
	ret := _m.Called(ctx, ns, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, ns, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Events provides a mock function with given fields:
func (_m *Orchestrator) Events() events.EventManager {
	ret := _m.Called()

	var r0 events.EventManager
	if rf, ok := ret.Get(0).(func() events.EventManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(events.EventManager)
		}
	}

	return r0
}

// GetBatchByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetBatchByID(ctx context.Context, ns string, id string) (*fftypes.Batch, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Batch
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Batch); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Batch)
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

// GetBatches provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetBatches(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Batch, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Batch
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Batch); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Batch)
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

// GetChartHistogram provides a mock function with given fields: ctx, ns, startTime, endTime, buckets, tableName
func (_m *Orchestrator) GetChartHistogram(ctx context.Context, ns string, startTime int64, endTime int64, buckets int64, tableName database.CollectionName) ([]*fftypes.ChartHistogram, error) {
	ret := _m.Called(ctx, ns, startTime, endTime, buckets, tableName)

	var r0 []*fftypes.ChartHistogram
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64, int64, database.CollectionName) []*fftypes.ChartHistogram); ok {
		r0 = rf(ctx, ns, startTime, endTime, buckets, tableName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.ChartHistogram)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int64, int64, int64, database.CollectionName) error); ok {
		r1 = rf(ctx, ns, startTime, endTime, buckets, tableName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfig provides a mock function with given fields: ctx
func (_m *Orchestrator) GetConfig(ctx context.Context) fftypes.JSONObject {
	ret := _m.Called(ctx)

	var r0 fftypes.JSONObject
	if rf, ok := ret.Get(0).(func(context.Context) fftypes.JSONObject); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fftypes.JSONObject)
		}
	}

	return r0
}

// GetConfigRecord provides a mock function with given fields: ctx, key
func (_m *Orchestrator) GetConfigRecord(ctx context.Context, key string) (*fftypes.ConfigRecord, error) {
	ret := _m.Called(ctx, key)

	var r0 *fftypes.ConfigRecord
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.ConfigRecord); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.ConfigRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfigRecords provides a mock function with given fields: ctx, filter
func (_m *Orchestrator) GetConfigRecords(ctx context.Context, filter database.AndFilter) ([]*fftypes.ConfigRecord, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.ConfigRecord
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.ConfigRecord); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.ConfigRecord)
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

// GetData provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetData(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Data, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Data
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Data); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Data)
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

// GetDataByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetDataByID(ctx context.Context, ns string, id string) (*fftypes.Data, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Data
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Data); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Data)
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

// GetDatatypeByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetDatatypeByID(ctx context.Context, ns string, id string) (*fftypes.Datatype, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Datatype
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Datatype); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Datatype)
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

// GetDatatypeByName provides a mock function with given fields: ctx, ns, name, version
func (_m *Orchestrator) GetDatatypeByName(ctx context.Context, ns string, name string, version string) (*fftypes.Datatype, error) {
	ret := _m.Called(ctx, ns, name, version)

	var r0 *fftypes.Datatype
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *fftypes.Datatype); ok {
		r0 = rf(ctx, ns, name, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Datatype)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, ns, name, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDatatypes provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetDatatypes(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Datatype, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Datatype
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Datatype); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Datatype)
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

// GetEventByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetEventByID(ctx context.Context, ns string, id string) (*fftypes.Event, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Event
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Event); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Event)
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

// GetEvents provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetEvents(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Event, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Event
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Event); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Event)
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

// GetMessageByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetMessageByID(ctx context.Context, ns string, id string) (*fftypes.Message, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Message); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Message)
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

// GetMessageByIDWithData provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetMessageByIDWithData(ctx context.Context, ns string, id string) (*fftypes.MessageInOut, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.MessageInOut
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.MessageInOut); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.MessageInOut)
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

// GetMessageData provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetMessageData(ctx context.Context, ns string, id string) ([]*fftypes.Data, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 []*fftypes.Data
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*fftypes.Data); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Data)
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

// GetMessageEvents provides a mock function with given fields: ctx, ns, id, filter
func (_m *Orchestrator) GetMessageEvents(ctx context.Context, ns string, id string, filter database.AndFilter) ([]*fftypes.Event, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id, filter)

	var r0 []*fftypes.Event
	if rf, ok := ret.Get(0).(func(context.Context, string, string, database.AndFilter) []*fftypes.Event); ok {
		r0 = rf(ctx, ns, id, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Event)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, id, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, id, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMessageOperations provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetMessageOperations(ctx context.Context, ns string, id string) ([]*fftypes.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 []*fftypes.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*fftypes.Operation); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Operation)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *database.FilterResult); ok {
		r1 = rf(ctx, ns, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, ns, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMessageTransaction provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetMessageTransaction(ctx context.Context, ns string, id string) (*fftypes.Transaction, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Transaction); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Transaction)
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

// GetMessages provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetMessages(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Message, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Message); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Message)
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

// GetMessagesForData provides a mock function with given fields: ctx, ns, dataID, filter
func (_m *Orchestrator) GetMessagesForData(ctx context.Context, ns string, dataID string, filter database.AndFilter) ([]*fftypes.Message, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, dataID, filter)

	var r0 []*fftypes.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, string, database.AndFilter) []*fftypes.Message); ok {
		r0 = rf(ctx, ns, dataID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Message)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, string, database.AndFilter) *database.FilterResult); ok {
		r1 = rf(ctx, ns, dataID, filter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string, database.AndFilter) error); ok {
		r2 = rf(ctx, ns, dataID, filter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetMessagesWithData provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetMessagesWithData(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.MessageInOut, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.MessageInOut
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.MessageInOut); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.MessageInOut)
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

// GetNamespace provides a mock function with given fields: ctx, ns
func (_m *Orchestrator) GetNamespace(ctx context.Context, ns string) (*fftypes.Namespace, error) {
	ret := _m.Called(ctx, ns)

	var r0 *fftypes.Namespace
	if rf, ok := ret.Get(0).(func(context.Context, string) *fftypes.Namespace); ok {
		r0 = rf(ctx, ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Namespace)
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

// GetNamespaces provides a mock function with given fields: ctx, filter
func (_m *Orchestrator) GetNamespaces(ctx context.Context, filter database.AndFilter) ([]*fftypes.Namespace, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*fftypes.Namespace
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*fftypes.Namespace); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Namespace)
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

// GetOperationByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetOperationByID(ctx context.Context, ns string, id string) (*fftypes.Operation, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Operation); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Operation)
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

// GetOperations provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetOperations(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Operation); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Operation)
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

// GetStatus provides a mock function with given fields: ctx
func (_m *Orchestrator) GetStatus(ctx context.Context) (*fftypes.NodeStatus, error) {
	ret := _m.Called(ctx)

	var r0 *fftypes.NodeStatus
	if rf, ok := ret.Get(0).(func(context.Context) *fftypes.NodeStatus); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.NodeStatus)
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

// GetSubscriptionByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetSubscriptionByID(ctx context.Context, ns string, id string) (*fftypes.Subscription, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Subscription); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Subscription)
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

// GetSubscriptions provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetSubscriptions(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Subscription, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Subscription); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Subscription)
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

// GetTransactionByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetTransactionByID(ctx context.Context, ns string, id string) (*fftypes.Transaction, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *fftypes.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *fftypes.Transaction); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Transaction)
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

// GetTransactionOperations provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetTransactionOperations(ctx context.Context, ns string, id string) ([]*fftypes.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 []*fftypes.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*fftypes.Operation); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Operation)
		}
	}

	var r1 *database.FilterResult
	if rf, ok := ret.Get(1).(func(context.Context, string, string) *database.FilterResult); ok {
		r1 = rf(ctx, ns, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*database.FilterResult)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, ns, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTransactions provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetTransactions(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Transaction, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*fftypes.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*fftypes.Transaction); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*fftypes.Transaction)
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

// Init provides a mock function with given fields: ctx, cancelCtx
func (_m *Orchestrator) Init(ctx context.Context, cancelCtx context.CancelFunc) error {
	ret := _m.Called(ctx, cancelCtx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, context.CancelFunc) error); ok {
		r0 = rf(ctx, cancelCtx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsPreInit provides a mock function with given fields:
func (_m *Orchestrator) IsPreInit() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NetworkMap provides a mock function with given fields:
func (_m *Orchestrator) NetworkMap() networkmap.Manager {
	ret := _m.Called()

	var r0 networkmap.Manager
	if rf, ok := ret.Get(0).(func() networkmap.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(networkmap.Manager)
		}
	}

	return r0
}

// PrivateMessaging provides a mock function with given fields:
func (_m *Orchestrator) PrivateMessaging() privatemessaging.Manager {
	ret := _m.Called()

	var r0 privatemessaging.Manager
	if rf, ok := ret.Get(0).(func() privatemessaging.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(privatemessaging.Manager)
		}
	}

	return r0
}

// PutConfigRecord provides a mock function with given fields: ctx, key, configRecord
func (_m *Orchestrator) PutConfigRecord(ctx context.Context, key string, configRecord *fftypes.JSONAny) (*fftypes.JSONAny, error) {
	ret := _m.Called(ctx, key, configRecord)

	var r0 *fftypes.JSONAny
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.JSONAny) *fftypes.JSONAny); ok {
		r0 = rf(ctx, key, configRecord)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.JSONAny)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.JSONAny) error); ok {
		r1 = rf(ctx, key, configRecord)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RequestReply provides a mock function with given fields: ctx, ns, msg
func (_m *Orchestrator) RequestReply(ctx context.Context, ns string, msg *fftypes.MessageInOut) (*fftypes.MessageInOut, error) {
	ret := _m.Called(ctx, ns, msg)

	var r0 *fftypes.MessageInOut
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.MessageInOut) *fftypes.MessageInOut); ok {
		r0 = rf(ctx, ns, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.MessageInOut)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.MessageInOut) error); ok {
		r1 = rf(ctx, ns, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetConfig provides a mock function with given fields: ctx
func (_m *Orchestrator) ResetConfig(ctx context.Context) {
	_m.Called(ctx)
}

// Start provides a mock function with given fields:
func (_m *Orchestrator) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WaitStop provides a mock function with given fields:
func (_m *Orchestrator) WaitStop() {
	_m.Called()
}
