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

package broadcast

import (
	"context"
	"fmt"
	"testing"

	"github.com/hyperledger-labs/firefly/internal/config"
	"github.com/hyperledger-labs/firefly/mocks/batchmocks"
	"github.com/hyperledger-labs/firefly/mocks/batchpinmocks"
	"github.com/hyperledger-labs/firefly/mocks/blockchainmocks"
	"github.com/hyperledger-labs/firefly/mocks/databasemocks"
	"github.com/hyperledger-labs/firefly/mocks/dataexchangemocks"
	"github.com/hyperledger-labs/firefly/mocks/datamocks"
	"github.com/hyperledger-labs/firefly/mocks/identitymanagermocks"
	"github.com/hyperledger-labs/firefly/mocks/publicstoragemocks"
	"github.com/hyperledger-labs/firefly/mocks/syncasyncmocks"
	"github.com/hyperledger-labs/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestBroadcast(t *testing.T) (*broadcastManager, func()) {
	config.Reset()
	mdi := &databasemocks.Plugin{}
	mim := &identitymanagermocks.Manager{}
	mdm := &datamocks.Manager{}
	mbi := &blockchainmocks.Plugin{}
	mpi := &publicstoragemocks.Plugin{}
	mba := &batchmocks.Manager{}
	mdx := &dataexchangemocks.Plugin{}
	msa := &syncasyncmocks.Bridge{}
	mbp := &batchpinmocks.Submitter{}
	mbi.On("Name").Return("ut_blockchain").Maybe()
	mba.On("RegisterDispatcher", []fftypes.MessageType{fftypes.MessageTypeBroadcast, fftypes.MessageTypeDefinition}, mock.Anything, mock.Anything).Return()
	ctx, cancel := context.WithCancel(context.Background())
	b, err := NewBroadcastManager(ctx, mdi, mim, mdm, mbi, mdx, mpi, mba, msa, mbp)
	assert.NoError(t, err)
	return b.(*broadcastManager), cancel
}

func TestInitFail(t *testing.T) {
	_, err := NewBroadcastManager(context.Background(), nil, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.Regexp(t, "FF10128", err)
}

func TestBroadcastMessageGood(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	msg := &fftypes.Message{}
	bm.database.(*databasemocks.Plugin).On("InsertMessageLocal", mock.Anything, msg).Return(nil)

	msgRet, err := bm.broadcastMessageCommon(context.Background(), msg, false, false)
	assert.NoError(t, err)
	assert.Equal(t, msg, msgRet)

	bm.Start()
	bm.WaitStop()
}

func TestBroadcastMessageBad(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	dupID := fftypes.NewUUID()
	msg := &fftypes.Message{
		Data: fftypes.DataRefs{
			{ID: dupID /* missing hash */},
		},
	}
	bm.database.(*databasemocks.Plugin).On("UpsertMessage", mock.Anything, msg, false).Return(nil)

	_, err := bm.broadcastMessageCommon(context.Background(), msg, false, false)
	assert.Regexp(t, "FF10144", err)

}

func TestDispatchBatchInvalidData(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	err := bm.dispatchBatch(context.Background(), &fftypes.Batch{
		Payload: fftypes.BatchPayload{
			Data: []*fftypes.Data{
				{Value: fftypes.Byteable(`!json`)},
			},
		},
	}, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.Regexp(t, "FF10137", err)
}

func TestDispatchBatchUploadFail(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	bm.publicstorage.(*publicstoragemocks.Plugin).On("PublishData", mock.Anything, mock.Anything).Return("", fmt.Errorf("pop"))

	err := bm.dispatchBatch(context.Background(), &fftypes.Batch{}, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.EqualError(t, err, "pop")
}

func TestDispatchBatchSubmitBatchPinSucceed(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	mdi := bm.database.(*databasemocks.Plugin)
	mdi.On("RunAsGroup", mock.Anything, mock.Anything).Return(nil)
	bm.publicstorage.(*publicstoragemocks.Plugin).On("PublishData", mock.Anything, mock.Anything).Return("id1", nil)

	err := bm.dispatchBatch(context.Background(), &fftypes.Batch{}, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.NoError(t, err)
}

func TestDispatchBatchSubmitBroadcastFail(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	mdi := bm.database.(*databasemocks.Plugin)
	mps := bm.publicstorage.(*publicstoragemocks.Plugin)
	mbp := bm.batchpin.(*batchpinmocks.Submitter)
	mdi.On("RunAsGroup", mock.Anything, mock.Anything).Return(nil)
	mps.On("PublishData", mock.Anything, mock.Anything).Return("id1", nil)
	mps.On("Name").Return("ut_publicstorage")

	err := bm.dispatchBatch(context.Background(), &fftypes.Batch{Identity: fftypes.Identity{Author: "wrong", Key: "wrong"}}, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.NoError(t, err)

	mdi.On("UpdateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertOperation", mock.Anything, mock.Anything, false).Return(nil)
	mbp.On("SubmitPinnedBatch", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	fn := mdi.Calls[0].Arguments[1].(func(ctx context.Context) error)
	err = fn(context.Background())
	assert.Regexp(t, "pop", err)
}

func TestSubmitTXAndUpdateDBUpdateBatchFail(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	mdi := bm.database.(*databasemocks.Plugin)
	mdi.On("UpsertTransaction", mock.Anything, mock.Anything, false).Return(nil)
	mdi.On("UpdateBatch", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	bm.blockchain.(*blockchainmocks.Plugin).On("SubmitBatchPin", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", fmt.Errorf("pop"))

	err := bm.submitTXAndUpdateDB(context.Background(), &fftypes.Batch{Identity: fftypes.Identity{Author: "org1", Key: "0x12345"}}, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.Regexp(t, "pop", err)
}

func TestSubmitTXAndUpdateDBAddOp1Fail(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	mdi := bm.database.(*databasemocks.Plugin)
	mbi := bm.blockchain.(*blockchainmocks.Plugin)
	mdi.On("UpsertTransaction", mock.Anything, mock.Anything, false).Return(nil)
	mdi.On("UpdateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertOperation", mock.Anything, mock.Anything, false).Return(fmt.Errorf("pop"))
	mbi.On("SubmitBatchPin", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("txid", nil)
	mbi.On("Name").Return("unittest")
	bm.publicstorage.(*publicstoragemocks.Plugin).On("Name").Return("ut_publicstorage")

	batch := &fftypes.Batch{
		Identity: fftypes.Identity{Author: "org1", Key: "0x12345"},
		Payload: fftypes.BatchPayload{
			Messages: []*fftypes.Message{
				{Header: fftypes.MessageHeader{
					ID: fftypes.NewUUID(),
				}},
			},
		},
	}

	err := bm.submitTXAndUpdateDB(context.Background(), batch, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.Regexp(t, "pop", err)
}

func TestSubmitTXAndUpdateDBSucceed(t *testing.T) {
	bm, cancel := newTestBroadcast(t)
	defer cancel()

	mdi := bm.database.(*databasemocks.Plugin)
	mbi := bm.blockchain.(*blockchainmocks.Plugin)
	mbp := bm.batchpin.(*batchpinmocks.Submitter)
	mdi.On("UpsertTransaction", mock.Anything, mock.Anything, false).Return(nil)
	mdi.On("UpdateBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mdi.On("UpsertOperation", mock.Anything, mock.Anything, false).Return(nil)
	mbi.On("SubmitBatchPin", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mbp.On("SubmitPinnedBatch", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	bm.publicstorage.(*publicstoragemocks.Plugin).On("Name").Return("ut_publicstorage")

	msgID := fftypes.NewUUID()
	batch := &fftypes.Batch{
		Identity: fftypes.Identity{Author: "org1", Key: "0x12345"},
		Payload: fftypes.BatchPayload{
			TX: fftypes.TransactionRef{
				Type: fftypes.TransactionTypeBatchPin,
				ID:   fftypes.NewUUID(),
			},
			Messages: []*fftypes.Message{
				{Header: fftypes.MessageHeader{
					ID: msgID,
				}},
			},
		},
		PayloadRef: "ipfs_id",
	}

	err := bm.submitTXAndUpdateDB(context.Background(), batch, []*fftypes.Bytes32{fftypes.NewRandB32()})
	assert.NoError(t, err)

	op := mdi.Calls[1].Arguments[1].(*fftypes.Operation)
	assert.Equal(t, *batch.Payload.TX.ID, *op.Transaction)
	assert.Equal(t, "ut_publicstorage", op.Plugin)
	assert.Equal(t, "ipfs_id", op.BackendID)
	assert.Equal(t, fftypes.OpTypePublicStorageBatchBroadcast, op.Type)

}
