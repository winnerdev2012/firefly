// Copyright © 2022 Kaleido, Inc.
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

package orchestrator

import (
	"context"
	"fmt"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
)

func (or *orchestrator) getPlugins() core.NamespaceStatusPlugins {
	// Plugins can have more than one name, so they must be iterated over
	tokensArray := make([]*core.NamespaceStatusPlugin, 0)
	for _, plugin := range or.plugins.Tokens {
		tokensArray = append(tokensArray, &core.NamespaceStatusPlugin{
			Name:       plugin.Name,
			PluginType: plugin.Plugin.Name(),
		})
	}

	blockchainsArray := make([]*core.NamespaceStatusPlugin, 0)
	blockchainsArray = append(blockchainsArray, &core.NamespaceStatusPlugin{
		Name:       or.plugins.Blockchain.Name,
		PluginType: or.plugins.Blockchain.Plugin.Name(),
	})

	databasesArray := make([]*core.NamespaceStatusPlugin, 0)
	databasesArray = append(databasesArray, &core.NamespaceStatusPlugin{
		Name:       or.plugins.Database.Name,
		PluginType: or.plugins.Database.Plugin.Name(),
	})

	sharedstorageArray := make([]*core.NamespaceStatusPlugin, 0)
	sharedstorageArray = append(sharedstorageArray, &core.NamespaceStatusPlugin{
		Name:       or.plugins.SharedStorage.Name,
		PluginType: or.plugins.SharedStorage.Plugin.Name(),
	})

	dataexchangeArray := make([]*core.NamespaceStatusPlugin, 0)
	dataexchangeArray = append(dataexchangeArray, &core.NamespaceStatusPlugin{
		Name:       or.plugins.DataExchange.Name,
		PluginType: or.plugins.DataExchange.Plugin.Name(),
	})

	return core.NamespaceStatusPlugins{
		Blockchain:    blockchainsArray,
		Database:      databasesArray,
		SharedStorage: sharedstorageArray,
		DataExchange:  dataexchangeArray,
		Events:        or.events.GetPlugins(),
		Tokens:        tokensArray,
		Identity:      []*core.NamespaceStatusPlugin{},
	}
}

func (or *orchestrator) GetNodeUUID(ctx context.Context) (node *fftypes.UUID) {
	if or.node != nil {
		return or.node
	}
	status, err := or.GetStatus(ctx)
	if err != nil {
		log.L(or.ctx).Warnf("Failed to query local node UUID: %s", err)
		return nil
	}
	if status.Node.Registered {
		or.node = status.Node.ID
	} else {
		log.L(or.ctx).Infof("Node not yet registered")
	}
	return or.node
}

func (or *orchestrator) GetStatus(ctx context.Context) (status *core.NamespaceStatus, err error) {

	org, err := or.identity.GetMultipartyRootOrg(ctx)
	if err != nil {
		log.L(ctx).Warnf("Failed to query local org for status: %s", err)
	}

	status = &core.NamespaceStatus{
		Namespace: or.namespace,
		Node: core.NamespaceStatusNode{
			Name: config.GetString(coreconfig.NodeName),
		},
		Org: core.NamespaceStatusOrg{
			Name: or.config.Multiparty.Org.Name,
		},
		Plugins: or.getPlugins(),
		Multiparty: core.NamespaceStatusMultiparty{
			Enabled: or.config.Multiparty.Enabled,
		},
	}

	if or.config.Multiparty.Enabled {
		status.Multiparty.Contracts = or.namespace.Contracts
	}

	if org != nil {
		status.Org.Registered = true
		status.Org.ID = org.ID
		status.Org.DID = org.DID

		fb := database.VerifierQueryFactory.NewFilter(ctx)
		verifiers, _, err := or.database().GetVerifiers(ctx, org.Namespace, fb.And(fb.Eq("identity", org.ID)))
		if err != nil {
			return nil, err
		}
		status.Org.Verifiers = make([]*core.VerifierRef, len(verifiers))
		for i, v := range verifiers {
			status.Org.Verifiers[i] = &v.VerifierRef
		}

		node, _, err := or.identity.CachedIdentityLookupNilOK(ctx, fmt.Sprintf("%s%s", core.FireFlyNodeDIDPrefix, status.Node.Name))
		if err != nil {
			return nil, err
		}
		if node != nil && !node.Parent.Equals(org.ID) {
			log.L(ctx).Errorf("Specified node name is in use by another org: %s", err)
			node = nil
		}
		if node != nil {
			status.Node.Registered = true
			status.Node.ID = node.ID
		}
	}

	return status, nil
}
