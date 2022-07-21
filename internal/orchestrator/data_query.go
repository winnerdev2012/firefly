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
	"database/sql/driver"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
)

func (or *orchestrator) GetNamespace(ctx context.Context) *core.Namespace {
	return or.namespace
}

func (or *orchestrator) GetTransactionByID(ctx context.Context, id string) (*core.Transaction, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetTransactionByID(ctx, or.namespace.LocalName, u)
}

func (or *orchestrator) GetTransactionOperations(ctx context.Context, id string) ([]*core.Operation, *database.FilterResult, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	fb := database.OperationQueryFactory.NewFilter(ctx)
	filter := fb.And(
		fb.Eq("tx", u),
	)
	return or.database().GetOperations(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) getMessageByID(ctx context.Context, id string) (*core.Message, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	msg, err := or.database().GetMessageByID(ctx, or.namespace.LocalName, u)
	if err == nil && msg == nil {
		return nil, i18n.NewError(ctx, coremsgs.Msg404NotFound)
	}
	return msg, err
}

func (or *orchestrator) GetMessageByID(ctx context.Context, id string) (*core.Message, error) {
	return or.getMessageByID(ctx, id)
}

func (or *orchestrator) fetchMessageData(ctx context.Context, msg *core.Message) (*core.MessageInOut, error) {
	msgI := &core.MessageInOut{
		Message: *msg,
	}
	// Lookup the full data
	data, _, err := or.data.GetMessageDataCached(ctx, msg)
	if err != nil {
		return nil, err
	}
	msgI.SetInlineData(data)
	return msgI, err
}

func (or *orchestrator) GetMessageByIDWithData(ctx context.Context, id string) (*core.MessageInOut, error) {
	msg, err := or.getMessageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.fetchMessageData(ctx, msg)
}

func (or *orchestrator) GetBatchByID(ctx context.Context, id string) (*core.BatchPersisted, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetBatchByID(ctx, or.namespace.LocalName, u)
}

func (or *orchestrator) GetDataByID(ctx context.Context, id string) (*core.Data, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetDataByID(ctx, or.namespace.LocalName, u, true)
}

func (or *orchestrator) GetDatatypeByID(ctx context.Context, id string) (*core.Datatype, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetDatatypeByID(ctx, or.namespace.LocalName, u)
}

func (or *orchestrator) GetDatatypeByName(ctx context.Context, name, version string) (*core.Datatype, error) {
	if err := fftypes.ValidateFFNameFieldNoUUID(ctx, name, "name"); err != nil {
		return nil, err
	}
	return or.database().GetDatatypeByName(ctx, or.namespace.LocalName, name, version)
}

func (or *orchestrator) GetOperationByID(ctx context.Context, id string) (*core.Operation, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetOperationByID(ctx, or.namespace.LocalName, u)
}

func (or *orchestrator) GetEventByID(ctx context.Context, id string) (*core.Event, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetEventByID(ctx, or.namespace.LocalName, u)
}

func (or *orchestrator) GetTransactions(ctx context.Context, filter database.AndFilter) ([]*core.Transaction, *database.FilterResult, error) {
	return or.database().GetTransactions(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetMessages(ctx context.Context, filter database.AndFilter) ([]*core.Message, *database.FilterResult, error) {
	return or.database().GetMessages(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetMessagesWithData(ctx context.Context, filter database.AndFilter) ([]*core.MessageInOut, *database.FilterResult, error) {
	msgs, fr, err := or.database().GetMessages(ctx, or.namespace.LocalName, filter)
	if err != nil {
		return nil, nil, err
	}
	msgsData := make([]*core.MessageInOut, len(msgs))
	for i, msg := range msgs {
		if msgsData[i], err = or.fetchMessageData(ctx, msg); err != nil {
			return nil, nil, err
		}
	}
	return msgsData, fr, err
}

func (or *orchestrator) GetMessageData(ctx context.Context, id string) (core.DataArray, error) {
	msg, err := or.getMessageByID(ctx, id)
	if err != nil || msg == nil {
		return nil, err
	}
	data, _, err := or.data.GetMessageDataCached(ctx, msg)
	return data, err
}

func (or *orchestrator) getMessageTransactionID(ctx context.Context, id string) (*fftypes.UUID, error) {
	msg, err := or.getMessageByID(ctx, id)
	if err != nil || msg == nil {
		return nil, err
	}
	var txID *fftypes.UUID
	if msg.Header.TxType == core.TransactionTypeBatchPin {
		if msg.BatchID == nil {
			return nil, i18n.NewError(ctx, coremsgs.MsgBatchNotSet)
		}
		batch, err := or.database().GetBatchByID(ctx, or.namespace.LocalName, msg.BatchID)
		if err != nil {
			return nil, err
		}
		if batch == nil {
			return nil, i18n.NewError(ctx, coremsgs.MsgBatchNotFound, msg.BatchID)
		}
		txID = batch.TX.ID
		if txID == nil {
			return nil, i18n.NewError(ctx, coremsgs.MsgBatchTXNotSet, msg.BatchID)
		}
	} else {
		return nil, i18n.NewError(ctx, coremsgs.MsgNoTransaction)
	}
	return txID, nil
}

func (or *orchestrator) GetMessageTransaction(ctx context.Context, id string) (*core.Transaction, error) {
	txID, err := or.getMessageTransactionID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetTransactionByID(ctx, or.namespace.LocalName, txID)
}

func (or *orchestrator) GetMessageEvents(ctx context.Context, id string, filter database.AndFilter) ([]*core.Event, *database.FilterResult, error) {
	msg, err := or.getMessageByID(ctx, id)
	if err != nil || msg == nil {
		return nil, nil, err
	}
	// Events can refer to the message, or any data in the message
	// So scope the event down to those referred UUIDs, in addition to any and conditions passed in
	referencedIDs := make([]driver.Value, len(msg.Data)+1)
	referencedIDs[0] = msg.Header.ID
	for i, dataRef := range msg.Data {
		referencedIDs[i+1] = dataRef.ID
	}
	filter = filter.Condition(filter.Builder().In("reference", referencedIDs))
	// Execute the filter
	return or.database().GetEvents(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetBatches(ctx context.Context, filter database.AndFilter) ([]*core.BatchPersisted, *database.FilterResult, error) {
	return or.database().GetBatches(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetData(ctx context.Context, filter database.AndFilter) (core.DataArray, *database.FilterResult, error) {
	return or.database().GetData(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetMessagesForData(ctx context.Context, id string, filter database.AndFilter) ([]*core.Message, *database.FilterResult, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return or.database().GetMessagesForData(ctx, or.namespace.LocalName, u, filter)
}

func (or *orchestrator) GetDatatypes(ctx context.Context, filter database.AndFilter) ([]*core.Datatype, *database.FilterResult, error) {
	return or.database().GetDatatypes(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetOperations(ctx context.Context, filter database.AndFilter) ([]*core.Operation, *database.FilterResult, error) {
	return or.database().GetOperations(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetEvents(ctx context.Context, filter database.AndFilter) ([]*core.Event, *database.FilterResult, error) {
	return or.database().GetEvents(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetBlockchainEventByID(ctx context.Context, id string) (*core.BlockchainEvent, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	return or.database().GetBlockchainEventByID(ctx, or.namespace.LocalName, u)
}

func (or *orchestrator) GetBlockchainEvents(ctx context.Context, filter database.AndFilter) ([]*core.BlockchainEvent, *database.FilterResult, error) {
	return or.database().GetBlockchainEvents(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetTransactionBlockchainEvents(ctx context.Context, id string) ([]*core.BlockchainEvent, *database.FilterResult, error) {
	u, err := fftypes.ParseUUID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	fb := database.BlockchainEventQueryFactory.NewFilter(ctx)
	return or.database().GetBlockchainEvents(ctx, or.namespace.LocalName, fb.And(fb.Eq("tx.id", u)))
}

func (or *orchestrator) GetPins(ctx context.Context, filter database.AndFilter) ([]*core.Pin, *database.FilterResult, error) {
	return or.database().GetPins(ctx, or.namespace.LocalName, filter)
}

func (or *orchestrator) GetEventsWithReferences(ctx context.Context, filter database.AndFilter) ([]*core.EnrichedEvent, *database.FilterResult, error) {
	events, fr, err := or.database().GetEvents(ctx, or.namespace.LocalName, filter)
	if err != nil {
		return nil, nil, err
	}

	enriched := make([]*core.EnrichedEvent, len(events))
	for i, event := range events {
		enrichedEvent, err := or.txHelper.EnrichEvent(or.ctx, event)
		if err != nil {
			return nil, nil, err
		}
		enriched[i] = enrichedEvent
	}
	return enriched, fr, err
}
