// Code generated by mockery v1.0.0. DO NOT EDIT.

package orchestratormocks

import (
	assets "github.com/hyperledger/firefly/internal/assets"
	batch "github.com/hyperledger/firefly/internal/batch"

	broadcast "github.com/hyperledger/firefly/internal/broadcast"

	context "context"

	contracts "github.com/hyperledger/firefly/internal/contracts"

	core "github.com/hyperledger/firefly/pkg/core"

	data "github.com/hyperledger/firefly/internal/data"

	database "github.com/hyperledger/firefly/pkg/database"

	events "github.com/hyperledger/firefly/internal/events"

	metrics "github.com/hyperledger/firefly/internal/metrics"

	mock "github.com/stretchr/testify/mock"

	networkmap "github.com/hyperledger/firefly/internal/networkmap"

	operations "github.com/hyperledger/firefly/internal/operations"

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

// BatchManager provides a mock function with given fields:
func (_m *Orchestrator) BatchManager() batch.Manager {
	ret := _m.Called()

	var r0 batch.Manager
	if rf, ok := ret.Get(0).(func() batch.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(batch.Manager)
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
func (_m *Orchestrator) CreateSubscription(ctx context.Context, ns string, subDef *core.Subscription) (*core.Subscription, error) {
	ret := _m.Called(ctx, ns, subDef)

	var r0 *core.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, *core.Subscription) *core.Subscription); ok {
		r0 = rf(ctx, ns, subDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *core.Subscription) error); ok {
		r1 = rf(ctx, ns, subDef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUpdateSubscription provides a mock function with given fields: ctx, ns, subDef
func (_m *Orchestrator) CreateUpdateSubscription(ctx context.Context, ns string, subDef *core.Subscription) (*core.Subscription, error) {
	ret := _m.Called(ctx, ns, subDef)

	var r0 *core.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, *core.Subscription) *core.Subscription); ok {
		r0 = rf(ctx, ns, subDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Subscription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *core.Subscription) error); ok {
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
func (_m *Orchestrator) GetBatchByID(ctx context.Context, ns string, id string) (*core.BatchPersisted, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.BatchPersisted
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.BatchPersisted); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.BatchPersisted)
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
func (_m *Orchestrator) GetBatches(ctx context.Context, ns string, filter database.AndFilter) ([]*core.BatchPersisted, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.BatchPersisted
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.BatchPersisted); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.BatchPersisted)
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

// GetBlockchainEventByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetBlockchainEventByID(ctx context.Context, ns string, id string) (*core.BlockchainEvent, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.BlockchainEvent
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.BlockchainEvent); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.BlockchainEvent)
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

// GetBlockchainEvents provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetBlockchainEvents(ctx context.Context, ns string, filter database.AndFilter) ([]*core.BlockchainEvent, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.BlockchainEvent
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.BlockchainEvent); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.BlockchainEvent)
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
func (_m *Orchestrator) GetChartHistogram(ctx context.Context, ns string, startTime int64, endTime int64, buckets int64, tableName database.CollectionName) ([]*core.ChartHistogram, error) {
	ret := _m.Called(ctx, ns, startTime, endTime, buckets, tableName)

	var r0 []*core.ChartHistogram
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64, int64, database.CollectionName) []*core.ChartHistogram); ok {
		r0 = rf(ctx, ns, startTime, endTime, buckets, tableName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.ChartHistogram)
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

// GetData provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetData(ctx context.Context, ns string, filter database.AndFilter) (core.DataArray, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 core.DataArray
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) core.DataArray); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.DataArray)
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
func (_m *Orchestrator) GetDataByID(ctx context.Context, ns string, id string) (*core.Data, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Data
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Data); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Data)
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
func (_m *Orchestrator) GetDatatypeByID(ctx context.Context, ns string, id string) (*core.Datatype, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Datatype
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Datatype); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Datatype)
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
func (_m *Orchestrator) GetDatatypeByName(ctx context.Context, ns string, name string, version string) (*core.Datatype, error) {
	ret := _m.Called(ctx, ns, name, version)

	var r0 *core.Datatype
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *core.Datatype); ok {
		r0 = rf(ctx, ns, name, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Datatype)
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
func (_m *Orchestrator) GetDatatypes(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Datatype, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Datatype
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Datatype); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Datatype)
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
func (_m *Orchestrator) GetEventByID(ctx context.Context, ns string, id string) (*core.Event, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Event
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Event); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Event)
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
func (_m *Orchestrator) GetEvents(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Event, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Event
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Event); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Event)
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

// GetEventsWithReferences provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetEventsWithReferences(ctx context.Context, ns string, filter database.AndFilter) ([]*core.EnrichedEvent, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.EnrichedEvent
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.EnrichedEvent); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.EnrichedEvent)
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
func (_m *Orchestrator) GetMessageByID(ctx context.Context, ns string, id string) (*core.Message, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Message); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Message)
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
func (_m *Orchestrator) GetMessageByIDWithData(ctx context.Context, ns string, id string) (*core.MessageInOut, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.MessageInOut
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.MessageInOut); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.MessageInOut)
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
func (_m *Orchestrator) GetMessageData(ctx context.Context, ns string, id string) (core.DataArray, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 core.DataArray
	if rf, ok := ret.Get(0).(func(context.Context, string, string) core.DataArray); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.DataArray)
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
func (_m *Orchestrator) GetMessageEvents(ctx context.Context, ns string, id string, filter database.AndFilter) ([]*core.Event, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id, filter)

	var r0 []*core.Event
	if rf, ok := ret.Get(0).(func(context.Context, string, string, database.AndFilter) []*core.Event); ok {
		r0 = rf(ctx, ns, id, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Event)
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
func (_m *Orchestrator) GetMessageOperations(ctx context.Context, ns string, id string) ([]*core.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 []*core.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*core.Operation); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Operation)
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
func (_m *Orchestrator) GetMessageTransaction(ctx context.Context, ns string, id string) (*core.Transaction, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Transaction); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Transaction)
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
func (_m *Orchestrator) GetMessages(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Message, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Message); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Message)
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
func (_m *Orchestrator) GetMessagesForData(ctx context.Context, ns string, dataID string, filter database.AndFilter) ([]*core.Message, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, dataID, filter)

	var r0 []*core.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, string, database.AndFilter) []*core.Message); ok {
		r0 = rf(ctx, ns, dataID, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Message)
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
func (_m *Orchestrator) GetMessagesWithData(ctx context.Context, ns string, filter database.AndFilter) ([]*core.MessageInOut, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.MessageInOut
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.MessageInOut); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.MessageInOut)
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
func (_m *Orchestrator) GetNamespace(ctx context.Context, ns string) (*core.Namespace, error) {
	ret := _m.Called(ctx, ns)

	var r0 *core.Namespace
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.Namespace); ok {
		r0 = rf(ctx, ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Namespace)
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

// GetOperationByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetOperationByID(ctx context.Context, ns string, id string) (*core.Operation, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Operation); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Operation)
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

// GetOperationByNamespacedID provides a mock function with given fields: ctx, nsOpID
func (_m *Orchestrator) GetOperationByNamespacedID(ctx context.Context, nsOpID string) (*core.Operation, error) {
	ret := _m.Called(ctx, nsOpID)

	var r0 *core.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.Operation); ok {
		r0 = rf(ctx, nsOpID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Operation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, nsOpID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperations provides a mock function with given fields: ctx, filter
func (_m *Orchestrator) GetOperations(ctx context.Context, filter database.AndFilter) ([]*core.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*core.Operation
	if rf, ok := ret.Get(0).(func(context.Context, database.AndFilter) []*core.Operation); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Operation)
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

// GetOperationsNamespaced provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetOperationsNamespaced(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Operation); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Operation)
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

// GetPins provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetPins(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Pin, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Pin
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Pin); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Pin)
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

// GetStatus provides a mock function with given fields: ctx, ns
func (_m *Orchestrator) GetStatus(ctx context.Context, ns string) (*core.NodeStatus, error) {
	ret := _m.Called(ctx, ns)

	var r0 *core.NodeStatus
	if rf, ok := ret.Get(0).(func(context.Context, string) *core.NodeStatus); ok {
		r0 = rf(ctx, ns)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.NodeStatus)
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

// GetSubscriptionByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetSubscriptionByID(ctx context.Context, ns string, id string) (*core.Subscription, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Subscription); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Subscription)
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
func (_m *Orchestrator) GetSubscriptions(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Subscription, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Subscription
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Subscription); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Subscription)
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

// GetTransactionBlockchainEvents provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetTransactionBlockchainEvents(ctx context.Context, ns string, id string) ([]*core.BlockchainEvent, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 []*core.BlockchainEvent
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*core.BlockchainEvent); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.BlockchainEvent)
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

// GetTransactionByID provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetTransactionByID(ctx context.Context, ns string, id string) (*core.Transaction, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.Transaction); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Transaction)
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
func (_m *Orchestrator) GetTransactionOperations(ctx context.Context, ns string, id string) ([]*core.Operation, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 []*core.Operation
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*core.Operation); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Operation)
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

// GetTransactionStatus provides a mock function with given fields: ctx, ns, id
func (_m *Orchestrator) GetTransactionStatus(ctx context.Context, ns string, id string) (*core.TransactionStatus, error) {
	ret := _m.Called(ctx, ns, id)

	var r0 *core.TransactionStatus
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *core.TransactionStatus); ok {
		r0 = rf(ctx, ns, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.TransactionStatus)
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

// GetTransactions provides a mock function with given fields: ctx, ns, filter
func (_m *Orchestrator) GetTransactions(ctx context.Context, ns string, filter database.AndFilter) ([]*core.Transaction, *database.FilterResult, error) {
	ret := _m.Called(ctx, ns, filter)

	var r0 []*core.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string, database.AndFilter) []*core.Transaction); ok {
		r0 = rf(ctx, ns, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*core.Transaction)
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

// Metrics provides a mock function with given fields:
func (_m *Orchestrator) Metrics() metrics.Manager {
	ret := _m.Called()

	var r0 metrics.Manager
	if rf, ok := ret.Get(0).(func() metrics.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metrics.Manager)
		}
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

// Operations provides a mock function with given fields:
func (_m *Orchestrator) Operations() operations.Manager {
	ret := _m.Called()

	var r0 operations.Manager
	if rf, ok := ret.Get(0).(func() operations.Manager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(operations.Manager)
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

// RequestReply provides a mock function with given fields: ctx, ns, msg
func (_m *Orchestrator) RequestReply(ctx context.Context, ns string, msg *core.MessageInOut) (*core.MessageInOut, error) {
	ret := _m.Called(ctx, ns, msg)

	var r0 *core.MessageInOut
	if rf, ok := ret.Get(0).(func(context.Context, string, *core.MessageInOut) *core.MessageInOut); ok {
		r0 = rf(ctx, ns, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.MessageInOut)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *core.MessageInOut) error); ok {
		r1 = rf(ctx, ns, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

// SubmitNetworkAction provides a mock function with given fields: ctx, ns, action
func (_m *Orchestrator) SubmitNetworkAction(ctx context.Context, ns string, action *core.NetworkAction) error {
	ret := _m.Called(ctx, ns, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *core.NetworkAction) error); ok {
		r0 = rf(ctx, ns, action)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WaitStop provides a mock function with given fields:
func (_m *Orchestrator) WaitStop() {
	_m.Called()
}
