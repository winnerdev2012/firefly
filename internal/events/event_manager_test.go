// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly/internal/cache"
	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/hyperledger/firefly/internal/events/system"
	"github.com/hyperledger/firefly/internal/txcommon"
	"github.com/hyperledger/firefly/mocks/assetmocks"
	"github.com/hyperledger/firefly/mocks/blockchainmocks"
	"github.com/hyperledger/firefly/mocks/broadcastmocks"
	"github.com/hyperledger/firefly/mocks/cachemocks"
	"github.com/hyperledger/firefly/mocks/databasemocks"
	"github.com/hyperledger/firefly/mocks/datamocks"
	"github.com/hyperledger/firefly/mocks/definitionsmocks"
	"github.com/hyperledger/firefly/mocks/eventsmocks"
	"github.com/hyperledger/firefly/mocks/identitymanagermocks"
	"github.com/hyperledger/firefly/mocks/metricsmocks"
	"github.com/hyperledger/firefly/mocks/multipartymocks"
	"github.com/hyperledger/firefly/mocks/operationmocks"
	"github.com/hyperledger/firefly/mocks/privatemessagingmocks"
	"github.com/hyperledger/firefly/mocks/shareddownloadmocks"
	"github.com/hyperledger/firefly/mocks/txcommonmocks"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
	"github.com/hyperledger/firefly/pkg/events"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testNodeID = fftypes.NewUUID()
var testNode = &core.Identity{
	IdentityBase: core.IdentityBase{
		ID: testNodeID,
	},
}

type testEventManager struct {
	eventManager

	cancel func()
	mdi    *databasemocks.Plugin
	mbi    *blockchainmocks.Plugin
	mim    *identitymanagermocks.Manager
	met    *eventsmocks.Plugin
	mdm    *datamocks.Manager
	msh    *definitionsmocks.Handler
	mds    *definitionsmocks.Sender
	mbm    *broadcastmocks.Manager
	mpm    *privatemessagingmocks.Manager
	mam    *assetmocks.Manager
	msd    *shareddownloadmocks.Manager
	mmi    *metricsmocks.Manager
	mom    *operationmocks.Manager
	mev    *eventsmocks.Plugin
	mmp    *multipartymocks.Manager
	mth    *txcommonmocks.Helper
}

func (tem *testEventManager) cleanup(t *testing.T) {
	tem.cancel()
	tem.mdi.AssertExpectations(t)
	tem.mbi.AssertExpectations(t)
	tem.mim.AssertExpectations(t)
	tem.met.AssertExpectations(t)
	tem.mdm.AssertExpectations(t)
	tem.msh.AssertExpectations(t)
	tem.mds.AssertExpectations(t)
	tem.mbm.AssertExpectations(t)
	tem.mpm.AssertExpectations(t)
	tem.mam.AssertExpectations(t)
	tem.msd.AssertExpectations(t)
	tem.mmi.AssertExpectations(t)
	tem.mom.AssertExpectations(t)
	tem.mev.AssertExpectations(t)
	tem.mmp.AssertExpectations(t)
	tem.mth.AssertExpectations(t)
}

func newTestEventManager(t *testing.T) *testEventManager {
	return newTestEventManagerCommon(t, false, false)
}

func newTestEventManagerWithMetrics(t *testing.T) *testEventManager {
	return newTestEventManagerCommon(t, true, false)
}

func newTestEventManagerWithDBConcurrency(t *testing.T) *testEventManager {
	return newTestEventManagerCommon(t, false, true)
}

func newTestEventManagerCommon(t *testing.T, metrics, dbconcurrency bool) *testEventManager {
	coreconfig.Reset()
	config.Set(coreconfig.BlobReceiverWorkerCount, 1)
	config.Set(coreconfig.BlobReceiverWorkerBatchTimeout, "1s")
	logrus.SetLevel(logrus.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	mdi := &databasemocks.Plugin{}
	mbi := &blockchainmocks.Plugin{}
	mim := &identitymanagermocks.Manager{}
	met := &eventsmocks.Plugin{}
	mdm := &datamocks.Manager{}
	msh := &definitionsmocks.Handler{}
	mds := &definitionsmocks.Sender{}
	mbm := &broadcastmocks.Manager{}
	mpm := &privatemessagingmocks.Manager{}
	mam := &assetmocks.Manager{}
	msd := &shareddownloadmocks.Manager{}
	mmi := &metricsmocks.Manager{}
	cmi := &cachemocks.Manager{}
	cmi.On("GetCache", mock.Anything).Return(cache.NewUmanagedCache(ctx, 100, 5*time.Minute), nil)
	mom := &operationmocks.Manager{}
	mev := &eventsmocks.Plugin{}
	events := map[string]events.Plugin{"websockets": mev}
	mmp := &multipartymocks.Manager{}
	txHelper := &txcommonmocks.Helper{}
	mmi.On("IsMetricsEnabled").Return(metrics).Maybe()
	if metrics {
		mmi.On("TransferConfirmed", mock.Anything).Maybe()
	}
	met.On("Name").Return("ut").Maybe()
	mbi.On("VerifierType").Return(core.VerifierTypeEthAddress).Maybe()
	mdi.On("Capabilities").Return(&database.Capabilities{Concurrency: dbconcurrency}).Maybe()
	mev.On("SetHandler", "ns1", mock.Anything).Return(nil).Maybe()
	mev.On("ValidateOptions", mock.Anything).Return(nil).Maybe()
	ns := &core.Namespace{Name: "ns1", NetworkName: "ns1"}
	emi, err := NewEventManager(ctx, ns, mdi, mbi, mim, msh, mdm, mds, mbm, mpm, mam, msd, mmi, mom, txHelper, events, mmp, cmi)
	em := emi.(*eventManager)
	mockRunAsGroupPassthrough(mdi)
	assert.NoError(t, err)
	return &testEventManager{
		eventManager: *em,
		cancel:       cancel,
		mdi:          mdi,
		mbi:          mbi,
		mim:          mim,
		met:          met,
		mdm:          mdm,
		msh:          msh,
		mds:          mds,
		mbm:          mbm,
		mpm:          mpm,
		mam:          mam,
		msd:          msd,
		mmi:          mmi,
		mom:          mom,
		mev:          mev,
		mmp:          mmp,
		mth:          txHelper,
	}
}

func mockRunAsGroupPassthrough(mdi *databasemocks.Plugin) {
	rag := mdi.On("RunAsGroup", mock.Anything, mock.Anything).Maybe()
	rag.RunFn = func(a mock.Arguments) {
		fn := a[1].(func(context.Context) error)
		rag.ReturnArguments = mock.Arguments{fn(a[0].(context.Context))}
	}
}

func TestStartStop(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	em.mdi.On("GetOffset", mock.Anything, core.OffsetTypeAggregator, aggregatorOffsetName).Return(&core.Offset{
		Type:    core.OffsetTypeAggregator,
		Name:    aggregatorOffsetName,
		Current: 12345,
		RowID:   333333,
	}, nil)
	em.mdi.On("GetPins", mock.Anything, "ns1", mock.Anything).Return([]*core.Pin{}, nil, nil)
	em.mdi.On("GetSubscriptions", mock.Anything, mock.Anything, mock.Anything).Return([]*core.Subscription{}, nil, nil)
	assert.NoError(t, em.Start())
	em.NewEvents() <- 12345
	em.NewPins() <- 12345
	em.cancel()
	em.WaitStop()
}

func TestStartStopBadDependencies(t *testing.T) {
	_, err := NewEventManager(context.Background(), &core.Namespace{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.Regexp(t, "FF10128", err)

}

func TestStartStopEventListenerFail(t *testing.T) {
	config.Set(coreconfig.EventTransportsEnabled, []string{"wrongun"})
	defer coreconfig.Reset()
	mdi := &databasemocks.Plugin{}
	mbi := &blockchainmocks.Plugin{}
	mim := &identitymanagermocks.Manager{}
	mdm := &datamocks.Manager{}
	msh := &definitionsmocks.Handler{}
	mds := &definitionsmocks.Sender{}
	mbm := &broadcastmocks.Manager{}
	mpm := &privatemessagingmocks.Manager{}
	mam := &assetmocks.Manager{}
	msd := &shareddownloadmocks.Manager{}
	mm := &metricsmocks.Manager{}
	mom := &operationmocks.Manager{}
	mev := &eventsmocks.Plugin{}
	events := map[string]events.Plugin{"websockets": mev}
	mmp := &multipartymocks.Manager{}
	ctx := context.Background()
	cmi := &cachemocks.Manager{}
	cmi.On("GetCache", mock.Anything).Return(cache.NewUmanagedCache(ctx, 100, 5*time.Minute), nil)
	txHelper, _ := txcommon.NewTransactionHelper(ctx, "ns1", mdi, mdm, cmi)
	mdi.On("Capabilities").Return(&database.Capabilities{Concurrency: false})
	mbi.On("VerifierType").Return(core.VerifierTypeEthAddress)
	mev.On("SetHandler", "ns1", mock.Anything).Return(fmt.Errorf("pop"))
	ns := &core.Namespace{Name: "ns1", NetworkName: "ns1"}
	_, err := NewEventManager(context.Background(), ns, mdi, mbi, mim, msh, mdm, mds, mbm, mpm, mam, msd, mm, mom, txHelper, events, mmp, cmi)
	assert.EqualError(t, err, "pop")
}

func TestEmitSubscriptionEventsNoops(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	em.mdi.On("GetOffset", mock.Anything, core.OffsetTypeAggregator, aggregatorOffsetName).Return(&core.Offset{
		Type:    core.OffsetTypeAggregator,
		Name:    aggregatorOffsetName,
		Current: 12345,
		RowID:   333333,
	}, nil)
	em.mdi.On("GetPins", mock.Anything, "ns1", mock.Anything).Return([]*core.Pin{}, nil, nil).Maybe()
	em.mdi.On("GetSubscriptions", mock.Anything, mock.Anything, mock.Anything).Return([]*core.Subscription{}, nil, nil)

	getSubCallReady := make(chan bool, 1)
	getSubCalled := make(chan bool)
	getSub := em.mdi.On("GetSubscriptionByID", mock.Anything, "ns1", mock.Anything).Return(nil, nil)
	getSub.RunFn = func(a mock.Arguments) {
		<-getSubCallReady
		getSubCalled <- true
	}

	delOffsetCalled := make(chan bool)
	delOffsetMock := em.mdi.On("DeleteOffset", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	delOffsetMock.RunFn = func(a mock.Arguments) {
		delOffsetCalled <- true
	}

	assert.NoError(t, em.Start())

	// Wait until the gets occur for these events, which will return nil
	getSubCallReady <- true
	em.NewSubscriptions() <- fftypes.NewUUID()
	em.SubscriptionUpdates() <- fftypes.NewUUID()
	<-getSubCalled

	em.DeletedSubscriptions() <- fftypes.NewUUID()
	close(getSubCallReady)
	<-delOffsetCalled
}

func TestCreateDurableSubscriptionBadSub(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	err := em.CreateUpdateDurableSubscription(em.ctx, &core.Subscription{}, false)
	assert.Regexp(t, "FF10189", err)
}

func TestCreateDurableSubscriptionDupName(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
	}
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(sub, nil)
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, true)
	assert.Regexp(t, "FF10193", err)
}

func TestCreateDurableSubscriptionDefaultSubCannotParse(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
		Filter: core.SubscriptionFilter{
			Events: "![[[[[",
		},
	}
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, true)
	assert.Regexp(t, "FF10171", err)
}

func TestCreateDurableSubscriptionBadFirstEvent(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	wrongFirstEvent := core.SubOptsFirstEvent("lobster")
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
		Options: core.SubscriptionOptions{
			SubscriptionCoreOptions: core.SubscriptionCoreOptions{
				FirstEvent: &wrongFirstEvent,
			},
		},
	}
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(nil, nil)
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, true)
	assert.Regexp(t, "FF10191", err)
}

func TestCreateDurableSubscriptionNegativeFirstEvent(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	wrongFirstEvent := core.SubOptsFirstEvent("-12345")
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
		Options: core.SubscriptionOptions{
			SubscriptionCoreOptions: core.SubscriptionCoreOptions{
				FirstEvent: &wrongFirstEvent,
			},
		},
	}
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(nil, nil)
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, true)
	assert.Regexp(t, "FF10192", err)
}

func TestCreateDurableSubscriptionGetHighestSequenceFailure(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
	}
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(nil, nil)
	em.mdi.On("GetEvents", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, fmt.Errorf("pop"))
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, true)
	assert.EqualError(t, err, "pop")
}

func TestCreateDurableSubscriptionOk(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
	}
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(nil, nil)
	em.mdi.On("GetEvents", mock.Anything, "ns1", mock.Anything).Return([]*core.Event{
		{Sequence: 12345},
	}, nil, nil)
	em.mdi.On("UpsertSubscription", mock.Anything, mock.Anything, false).Return(nil)
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, true)
	assert.NoError(t, err)
	// Check genreated fields
	assert.NotNil(t, sub.ID)
	assert.Equal(t, "websockets", sub.Transport)
	assert.Equal(t, "12345", string(*sub.Options.FirstEvent))
}

func TestUpdateDurableSubscriptionOk(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
	}
	var firstEvent core.SubOptsFirstEvent = "12345"
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(&core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID: fftypes.NewUUID(),
		},
		Options: core.SubscriptionOptions{
			SubscriptionCoreOptions: core.SubscriptionCoreOptions{
				FirstEvent: &firstEvent,
			},
		},
	}, nil) // return non-matching existing
	em.mdi.On("UpsertSubscription", mock.Anything, mock.Anything, true).Return(nil)
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, false)
	assert.NoError(t, err)
	// Check genreated fields
	assert.NotNil(t, sub.ID)
	assert.Equal(t, "websockets", sub.Transport)
	assert.Equal(t, "12345", string(*sub.Options.FirstEvent))
}

func TestUpdateDurableSubscriptionNoOp(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	no := false
	sub := &core.Subscription{
		SubscriptionRef: core.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Namespace: "ns1",
			Name:      "sub1",
		},
		Transport: "websockets",
		Options: core.SubscriptionOptions{
			SubscriptionCoreOptions: core.SubscriptionCoreOptions{
				WithData: &no,
			},
		},
	}
	var subExisting = *sub
	subExisting.Created = fftypes.Now()
	subExisting.Updated = fftypes.Now()
	subExisting.ID = fftypes.NewUUID()
	em.mdi.On("GetSubscriptionByName", mock.Anything, "ns1", "sub1").Return(&subExisting, nil) // return non-matching existing
	err := em.CreateUpdateDurableSubscription(em.ctx, sub, false)
	assert.NoError(t, err)
}

func TestCreateDeleteDurableSubscriptionOk(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	subId := fftypes.NewUUID()
	sub := &core.Subscription{SubscriptionRef: core.SubscriptionRef{ID: subId, Namespace: "ns1"}}
	em.mdi.On("DeleteSubscriptionByID", mock.Anything, "ns1", subId).Return(nil)
	err := em.DeleteDurableSubscription(em.ctx, sub)
	assert.NoError(t, err)
}

func TestAddInternalListener(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)
	ie := &system.Events{}
	cbs := &eventsmocks.Callbacks{}

	cbs.On("RegisterConnection", mock.Anything, mock.Anything).Return(nil)
	cbs.On("EphemeralSubscription", mock.Anything, "ns1", mock.Anything, mock.Anything).Return(nil)

	conf := config.RootSection("ut.events")
	ie.InitConfig(conf)
	ie.Init(em.ctx, conf)
	ie.SetHandler("ns1", cbs)
	em.internalEvents = ie
	err := em.AddSystemEventListener("ns1", func(event *core.EventDelivery) error { return nil })
	assert.NoError(t, err)

	cbs.AssertExpectations(t)
}

func TestGetPlugins(t *testing.T) {
	em := newTestEventManager(t)
	defer em.cleanup(t)

	expectedPlugins := []*core.NamespaceStatusPlugin{
		{
			PluginType: "websockets",
		},
	}

	assert.ElementsMatch(t, em.GetPlugins(), expectedPlugins)
}
