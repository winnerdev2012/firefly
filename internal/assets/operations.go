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

package assets

import (
	"context"
	"fmt"

	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/internal/txcommon"
	"github.com/hyperledger/firefly/pkg/fftypes"
	"github.com/hyperledger/firefly/pkg/i18n"
)

type createPoolData struct {
	Pool *fftypes.TokenPool `json:"pool"`
}

type activatePoolData struct {
	Pool           *fftypes.TokenPool `json:"pool"`
	BlockchainInfo fftypes.JSONObject `json:"blockchainInfo"`
}

type transferData struct {
	Pool     *fftypes.TokenPool     `json:"pool"`
	Transfer *fftypes.TokenTransfer `json:"transfer"`
}

type approvalData struct {
	Pool     *fftypes.TokenPool     `json:"pool"`
	Approval *fftypes.TokenApproval `json:"approval"`
}

func (am *assetManager) PrepareOperation(ctx context.Context, op *fftypes.Operation) (*fftypes.PreparedOperation, error) {
	switch op.Type {
	case fftypes.OpTypeTokenCreatePool:
		pool, err := txcommon.RetrieveTokenPoolCreateInputs(ctx, op)
		if err != nil {
			return nil, err
		}
		return opCreatePool(op, pool), nil

	case fftypes.OpTypeTokenActivatePool:
		poolID, blockchainInfo, err := txcommon.RetrieveTokenPoolActivateInputs(ctx, op)
		if err != nil {
			return nil, err
		}
		pool, err := am.database.GetTokenPoolByID(ctx, poolID)
		if err != nil {
			return nil, err
		} else if pool == nil {
			return nil, i18n.NewError(ctx, coremsgs.Msg404NotFound)
		}
		return opActivatePool(op, pool, blockchainInfo), nil

	case fftypes.OpTypeTokenTransfer:
		transfer, err := txcommon.RetrieveTokenTransferInputs(ctx, op)
		if err != nil {
			return nil, err
		}
		pool, err := am.database.GetTokenPoolByID(ctx, transfer.Pool)
		if err != nil {
			return nil, err
		} else if pool == nil {
			return nil, i18n.NewError(ctx, coremsgs.Msg404NotFound)
		}
		return opTransfer(op, pool, transfer), nil

	case fftypes.OpTypeTokenApproval:
		approval, err := txcommon.RetrieveTokenApprovalInputs(ctx, op)
		if err != nil {
			return nil, err
		}
		pool, err := am.database.GetTokenPoolByID(ctx, approval.Pool)
		if err != nil {
			return nil, err
		} else if pool == nil {
			return nil, i18n.NewError(ctx, coremsgs.Msg404NotFound)
		}
		return opApproval(op, pool, approval), nil

	default:
		return nil, i18n.NewError(ctx, coremsgs.MsgOperationNotSupported, op.Type)
	}
}

func (am *assetManager) RunOperation(ctx context.Context, op *fftypes.PreparedOperation) (outputs fftypes.JSONObject, complete bool, err error) {
	switch data := op.Data.(type) {
	case createPoolData:
		plugin, err := am.selectTokenPlugin(ctx, data.Pool.Connector)
		if err != nil {
			return nil, false, err
		}
		complete, err = plugin.CreateTokenPool(ctx, op.ID, data.Pool)
		return nil, complete, err

	case activatePoolData:
		plugin, err := am.selectTokenPlugin(ctx, data.Pool.Connector)
		if err != nil {
			return nil, false, err
		}
		complete, err = plugin.ActivateTokenPool(ctx, op.ID, data.Pool, data.BlockchainInfo)
		return nil, complete, err

	case transferData:
		plugin, err := am.selectTokenPlugin(ctx, data.Pool.Connector)
		if err != nil {
			return nil, false, err
		}
		switch data.Transfer.Type {
		case fftypes.TokenTransferTypeMint:
			return nil, false, plugin.MintTokens(ctx, op.ID, data.Pool.ProtocolID, data.Transfer)
		case fftypes.TokenTransferTypeTransfer:
			return nil, false, plugin.TransferTokens(ctx, op.ID, data.Pool.ProtocolID, data.Transfer)
		case fftypes.TokenTransferTypeBurn:
			return nil, false, plugin.BurnTokens(ctx, op.ID, data.Pool.ProtocolID, data.Transfer)
		default:
			panic(fmt.Sprintf("unknown transfer type: %v", data.Transfer.Type))
		}

	case approvalData:
		plugin, err := am.selectTokenPlugin(ctx, data.Pool.Connector)
		if err != nil {
			return nil, false, err
		}
		return nil, false, plugin.TokensApproval(ctx, op.ID, data.Pool.ProtocolID, data.Approval)

	default:
		return nil, false, i18n.NewError(ctx, coremsgs.MsgOperationDataIncorrect, op.Data)
	}
}

func opCreatePool(op *fftypes.Operation, pool *fftypes.TokenPool) *fftypes.PreparedOperation {
	return &fftypes.PreparedOperation{
		ID:   op.ID,
		Type: op.Type,
		Data: createPoolData{Pool: pool},
	}
}

func opActivatePool(op *fftypes.Operation, pool *fftypes.TokenPool, blockchainInfo fftypes.JSONObject) *fftypes.PreparedOperation {
	return &fftypes.PreparedOperation{
		ID:   op.ID,
		Type: op.Type,
		Data: activatePoolData{Pool: pool, BlockchainInfo: blockchainInfo},
	}
}

func opTransfer(op *fftypes.Operation, pool *fftypes.TokenPool, transfer *fftypes.TokenTransfer) *fftypes.PreparedOperation {
	return &fftypes.PreparedOperation{
		ID:   op.ID,
		Type: op.Type,
		Data: transferData{Pool: pool, Transfer: transfer},
	}
}

func opApproval(op *fftypes.Operation, pool *fftypes.TokenPool, approval *fftypes.TokenApproval) *fftypes.PreparedOperation {
	return &fftypes.PreparedOperation{
		ID:   op.ID,
		Type: op.Type,
		Data: approvalData{Pool: pool, Approval: approval},
	}
}
