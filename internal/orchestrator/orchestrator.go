// Copyright © 2021 Kaleido, Inc.
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

package orchestrator

import (
	"context"
	"fmt"

	"github.com/kaleido-io/firefly/internal/batch"
	"github.com/kaleido-io/firefly/internal/blockchain/bifactory"
	"github.com/kaleido-io/firefly/internal/broadcast"
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/database/difactory"
	"github.com/kaleido-io/firefly/internal/events"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/internal/publicstorage/psfactory"
	"github.com/kaleido-io/firefly/pkg/blockchain"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/kaleido-io/firefly/pkg/publicstorage"
)

var (
	blockchainConfig    = config.NewPluginConfig("blockchain")
	databaseConfig      = config.NewPluginConfig("database")
	publicstorageConfig = config.NewPluginConfig("publicstorage")
)

// Orchestrator is the main interface behind the API, implementing the actions
type Orchestrator interface {
	blockchain.Callbacks

	Init(ctx context.Context) error
	Start() error
	WaitStop() // The close itself is performed by canceling the context

	// Broadcasts
	BroadcastNamespace(ctx context.Context, s *fftypes.Namespace) (*fftypes.Message, error)
	BroadcastDatatype(ctx context.Context, ns string, s *fftypes.Datatype) (*fftypes.Message, error)

	// Data Query
	GetNamespace(ctx context.Context, ns string) (*fftypes.Namespace, error)
	GetNamespaces(ctx context.Context, filter database.AndFilter) ([]*fftypes.Namespace, error)
	GetTransactionByID(ctx context.Context, ns, id string) (*fftypes.Transaction, error)
	GetTransactions(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Transaction, error)
	GetMessageByID(ctx context.Context, ns, id string) (*fftypes.Message, error)
	GetMessages(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Message, error)
	GetMessageOperations(ctx context.Context, ns, id string, filter database.AndFilter) ([]*fftypes.Operation, error)
	GetMessageEvents(ctx context.Context, ns, id string, filter database.AndFilter) ([]*fftypes.Event, error)
	GetMessagesForData(ctx context.Context, ns, dataID string, filter database.AndFilter) ([]*fftypes.Message, error)
	GetBatchByID(ctx context.Context, ns, id string) (*fftypes.Batch, error)
	GetBatches(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Batch, error)
	GetDataByID(ctx context.Context, ns, id string) (*fftypes.Data, error)
	GetData(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Data, error)
	GetDatatypeByID(ctx context.Context, ns, id string) (*fftypes.Datatype, error)
	GetDatatypes(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Datatype, error)
	GetOperationByID(ctx context.Context, ns, id string) (*fftypes.Operation, error)
	GetOperations(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Operation, error)
	GetEventByID(ctx context.Context, ns, id string) (*fftypes.Event, error)
	GetEvents(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Event, error)
}

type orchestrator struct {
	ctx           context.Context
	started       bool
	database      database.Plugin
	blockchain    blockchain.Plugin
	publicstorage publicstorage.Plugin
	events        events.EventManager
	batch         batch.Manager
	broadcast     broadcast.Manager
	nodeIDentity  string
}

func NewOrchestrator() Orchestrator {
	or := &orchestrator{}

	// Initialize the config on all the factories
	bifactory.InitPrefix(blockchainConfig)
	difactory.InitPrefix(databaseConfig)
	psfactory.InitPrefix(publicstorageConfig)

	return or
}

func (or *orchestrator) Init(ctx context.Context) (err error) {
	or.ctx = ctx
	err = or.initPlugins(ctx)
	if err == nil {
		err = or.initComponents(ctx)
	}
	if err == nil {
		err = or.initNamespaces(ctx)
	}
	return err
}

func (or *orchestrator) Start() error {
	err := or.blockchain.Start()
	if err == nil {
		err = or.batch.Start()
	}
	if err == nil {
		err = or.events.Start()
	}
	if err == nil {
		err = or.broadcast.Start()
	}
	or.started = true
	return err
}

func (or *orchestrator) WaitStop() {
	if !or.started {
		return
	}
	if or.batch != nil {
		or.batch.WaitStop()
		or.batch = nil
	}
	if or.broadcast != nil {
		or.broadcast.WaitStop()
		or.broadcast = nil
	}
	or.started = false
}

func (or *orchestrator) initPlugins(ctx context.Context) (err error) {

	if or.database == nil {
		diType := config.GetString(config.DatabaseType)
		if or.database, err = difactory.GetPlugin(ctx, diType); err != nil {
			return err
		}
	}
	if err = or.database.Init(ctx, databaseConfig.SubPrefix(or.database.Name()), or); err != nil {
		return err
	}

	if or.blockchain == nil {
		biType := config.GetString(config.BlockchainType)
		if or.blockchain, err = bifactory.GetPlugin(ctx, biType); err != nil {
			return err
		}
	}
	if err = or.initBlockchainPlugin(ctx); err != nil {
		return err
	}

	if or.publicstorage == nil {
		psType := config.GetString(config.PublicStorageType)
		if or.publicstorage, err = psfactory.GetPlugin(ctx, psType); err != nil {
			return err
		}
	}
	return or.publicstorage.Init(ctx, publicstorageConfig.SubPrefix(or.publicstorage.Name()), or)
}

func (or *orchestrator) initComponents(ctx context.Context) (err error) {
	if or.events == nil {
		or.events, err = events.NewEventManager(ctx, or.publicstorage, or.database)
		if err != nil {
			return err
		}
	}

	if or.batch == nil {
		or.batch, err = batch.NewBatchManager(ctx, or.database)
		if err != nil {
			return err
		}
	}

	if or.broadcast == nil {
		if or.broadcast, err = broadcast.NewBroadcastManager(ctx, or.database, or.blockchain, or.publicstorage, or.batch); err != nil {
			return err
		}
	}
	return nil
}

func (or *orchestrator) initBlockchainPlugin(ctx context.Context) error {
	err := or.blockchain.Init(ctx, blockchainConfig.SubPrefix(or.blockchain.Name()), or)
	if err != nil {
		return err
	}
	suppliedIDentity := config.GetString(config.NodeIdentity)
	or.nodeIDentity, err = or.blockchain.VerifyIdentitySyntax(ctx, suppliedIDentity)
	if err != nil {
		log.L(ctx).Errorf("Invalid node identity: %s", suppliedIDentity)
		return err
	}
	return nil
}

func (or *orchestrator) initNamespaces(ctx context.Context) error {
	defaultNS := config.GetString(config.NamespacesDefault)
	predefined := config.GetObjectArray(config.NamespacesPredefined)
	foundDefault := false
	for i, nsObject := range predefined {
		name := nsObject.GetString(ctx, "name")
		description := nsObject.GetString(ctx, "description")
		err := fftypes.ValidateFFNameField(ctx, name, fmt.Sprintf("namespaces.predefined[%d].name", i))
		if err != nil {
			return err
		}
		foundDefault = foundDefault || name == defaultNS

		ns, err := or.database.GetNamespace(ctx, name)
		if err != nil {
			return err
		}
		var updated bool
		if ns == nil {
			updated = true
			ns = &fftypes.Namespace{
				ID:          fftypes.NewUUID(),
				Name:        name,
				Description: description,
				Created:     fftypes.Now(),
			}
		} else {
			// Only update if the description has changed, and the one in our DB is locally defined
			updated = ns.Description != description && ns.Type == fftypes.NamespaceTypeLocal
			ns.Description = description
		}
		if updated {
			if err := or.database.UpsertNamespace(ctx, ns, true); err != nil {
				return err
			}
		}
	}
	if !foundDefault {
		return i18n.NewError(ctx, i18n.MsgDefaultNamespaceNotFound, defaultNS)
	}
	return nil
}
