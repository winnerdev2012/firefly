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

package multiparty

import (
	"context"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/internal/operations"
	"github.com/hyperledger/firefly/pkg/blockchain"
	"github.com/hyperledger/firefly/pkg/core"
)

type batchPinData struct {
	Batch      *core.BatchPersisted `json:"batch"`
	Contexts   []*fftypes.Bytes32   `json:"contexts"`
	PayloadRef string               `json:"payloadRef"`
}

type networkActionData struct {
	Type core.NetworkActionType `json:"type"`
	Key  string                 `json:"key"`
}

func addBatchPinInputs(op *core.Operation, batchID *fftypes.UUID, contexts []*fftypes.Bytes32, payloadRef string) {
	contextStr := make([]string, len(contexts))
	for i, c := range contexts {
		contextStr[i] = c.String()
	}
	op.Input = fftypes.JSONObject{
		"batch":      batchID.String(),
		"contexts":   contextStr,
		"payloadRef": payloadRef,
	}
}

func addNetworkActionInputs(op *core.Operation, actionType core.NetworkActionType, signingKey string) {
	op.Input = fftypes.JSONObject{
		"type": actionType.String(),
		"key":  signingKey,
	}
}

func retrieveBatchPinInputs(ctx context.Context, op *core.Operation) (batchID *fftypes.UUID, contexts []*fftypes.Bytes32, payloadRef string, err error) {
	batchID, err = fftypes.ParseUUID(ctx, op.Input.GetString("batch"))
	if err != nil {
		return nil, nil, "", err
	}
	contextStr := op.Input.GetStringArray("contexts")
	contexts = make([]*fftypes.Bytes32, len(contextStr))
	for i, c := range contextStr {
		contexts[i], err = fftypes.ParseBytes32(ctx, c)
		if err != nil {
			return nil, nil, "", err
		}
	}
	payloadRef = op.Input.GetString("payloadRef")
	return batchID, contexts, payloadRef, nil
}

func retrieveNetworkActionInputs(op *core.Operation) (actionType core.NetworkActionType, signingKey string) {
	actionType = fftypes.FFEnum(op.Input.GetString("type"))
	signingKey = op.Input.GetString("key")
	return actionType, signingKey
}

func (mm *multipartyManager) PrepareOperation(ctx context.Context, op *core.Operation) (*core.PreparedOperation, error) {
	switch op.Type {
	case core.OpTypeBlockchainPinBatch:
		batchID, contexts, payloadRef, err := retrieveBatchPinInputs(ctx, op)
		if err != nil {
			return nil, err
		}
		batch, err := mm.database.GetBatchByID(ctx, mm.namespace.LocalName, batchID)
		if err != nil {
			return nil, err
		} else if batch == nil {
			return nil, i18n.NewError(ctx, coremsgs.Msg404NotFound)
		}
		return opBatchPin(op, batch, contexts, payloadRef), nil

	case core.OpTypeBlockchainNetworkAction:
		actionType, signingKey := retrieveNetworkActionInputs(op)
		return opNetworkAction(op, actionType, signingKey), nil

	default:
		return nil, i18n.NewError(ctx, coremsgs.MsgOperationNotSupported, op.Type)
	}
}

func (mm *multipartyManager) RunOperation(ctx context.Context, op *core.PreparedOperation) (outputs fftypes.JSONObject, complete bool, err error) {
	switch data := op.Data.(type) {
	case batchPinData:
		batch := data.Batch

		// Only include namespace for V1 networks
		var namespace string
		if mm.activeContract.networkVersion == 1 {
			namespace = mm.namespace.RemoteName
		}

		return nil, false, mm.blockchain.SubmitBatchPin(ctx, op.NamespacedIDString(), batch.Key, &blockchain.BatchPin{
			Namespace:       namespace,
			TransactionID:   batch.TX.ID,
			BatchID:         batch.ID,
			BatchHash:       batch.Hash,
			BatchPayloadRef: data.PayloadRef,
			Contexts:        data.Contexts,
		}, mm.activeContract.location)

	case networkActionData:
		return nil, false, mm.blockchain.SubmitNetworkAction(ctx, op.NamespacedIDString(), data.Key, data.Type, mm.activeContract.location)

	default:
		return nil, false, i18n.NewError(ctx, coremsgs.MsgOperationDataIncorrect, op.Data)
	}
}

func (mm *multipartyManager) OnOperationUpdate(ctx context.Context, op *core.Operation, update *operations.OperationUpdate) error {
	return nil
}

func opBatchPin(op *core.Operation, batch *core.BatchPersisted, contexts []*fftypes.Bytes32, payloadRef string) *core.PreparedOperation {
	return &core.PreparedOperation{
		ID:        op.ID,
		Namespace: op.Namespace,
		Type:      op.Type,
		Data:      batchPinData{Batch: batch, Contexts: contexts, PayloadRef: payloadRef},
	}
}

func opNetworkAction(op *core.Operation, actionType core.NetworkActionType, key string) *core.PreparedOperation {
	return &core.PreparedOperation{
		ID:        op.ID,
		Namespace: op.Namespace,
		Type:      op.Type,
		Data: networkActionData{
			Type: actionType,
			Key:  key,
		},
	}
}
