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
	"encoding/json"
	"fmt"

	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/dataexchange"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

func (em *eventManager) MessageReceived(dx dataexchange.Plugin, peerID string, data []byte) {
	l := log.L(em.ctx)
	l.Infof("Message received from '%s' (len=%d)", peerID, len(data))

	// Try to de-serialize it as a batch
	var batch fftypes.Batch
	err := json.Unmarshal(data, &batch)
	if err != nil {
		l.Errorf("Invalid batch: %s", err)
		return
	}

	// Find the node associated with the peer
	filter := database.NodeQueryFactory.NewFilter(em.ctx).Eq("dx.peer", peerID)
	nodes, err := em.database.GetNodes(em.ctx, filter)
	if err != nil || len(nodes) < 1 {
		l.Errorf("Failed to retrieve node: %v", err)
		return
	}
	node := nodes[0]

	// Find the identity in the mesage
	batchOrg, err := em.database.GetOrganizationByIdentity(em.ctx, batch.Author)
	if err != nil || batchOrg == nil {
		l.Errorf("Failed to retrieve batch org: %v", err)
		return
	}

	// One of the orgs in the hierarchy of the batch author must be the owner of this node
	parent := batchOrg.Identity
	foundNodeOrg := batch.Author == node.Owner
	for !foundNodeOrg && parent != "" {
		candidate, err := em.database.GetOrganizationByIdentity(em.ctx, parent)
		if err != nil || candidate == nil {
			l.Errorf("Failed to retrieve node org '%s': %v", parent, err)
			return
		}
		foundNodeOrg = candidate.Identity == node.Owner
		parent = candidate.Parent
	}
	if !foundNodeOrg {
		l.Errorf("No org in the chain matches owner '%s' of node '%s' ('%s')", node.Owner, node.ID, node.Name)
		return
	}

	if err := em.persistBatch(em.ctx, &batch); err != nil {
		l.Errorf("Batch received from %s/%s invalid: %s", node.Owner, node.Name, err)
		return
	}

	em.aggregator.offchainBatches <- batch.ID
}

func (em *eventManager) BLOBReceived(dx dataexchange.Plugin, peerID string, ns string, id fftypes.UUID) {
}

func (em *eventManager) TransferResult(dx dataexchange.Plugin, trackingID string, status fftypes.OpStatus, info string, additionalInfo fftypes.JSONObject) {
	log.L(em.ctx).Infof("Transfer result %s=%s info='%s'", trackingID, status, info)

	// Find a matching operation, for this plugin, with the specified ID.
	// We retry a few times, as there's an outside possibility of the event arriving before we're finished persisting the operation itself
	var operations []*fftypes.Operation
	fb := database.OperationQueryFactory.NewFilter(em.ctx)
	filter := fb.And(
		fb.Eq("backendid", trackingID),
		fb.Eq("plugin", dx.Name()),
	)
	err := em.retry.Do(em.ctx, fmt.Sprintf("correlate transfer %s", trackingID), func(attempt int) (retry bool, err error) {
		operations, err = em.database.GetOperations(em.ctx, filter)
		if err == nil && len(operations) == 0 {
			err = i18n.NewError(em.ctx, i18n.Msg404NotFound)
		}
		return (err != nil && attempt <= em.opCorrelationRetries), err
	})
	if err != nil {
		log.L(em.ctx).Warnf("Failed to correlate transfer ID '%s' with a submitted operation", trackingID)
		return
	}

	update := database.OperationQueryFactory.NewUpdate(em.ctx).
		Set("status", status).
		Set("error", info).
		Set("info", additionalInfo)
	for _, op := range operations {
		if err := em.database.UpdateOperation(em.ctx, op.ID, update); err != nil {
			return
		}
	}

}
