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

package defsender

import (
	"context"
	"encoding/json"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/broadcast"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/internal/data"
	"github.com/hyperledger/firefly/internal/identity"
	"github.com/hyperledger/firefly/pkg/core"
)

type DefinitionHandler interface {
	HandleDefinition(ctx context.Context, state *core.BatchState, msg *core.Message, data *core.Data) error
}

type Sender interface {
	core.Named

	Init(handler DefinitionHandler)
	CreateDefinition(ctx context.Context, def core.Definition, tag string, waitConfirm bool) (msg *core.Message, err error)
	CreateDefinitionWithIdentity(ctx context.Context, def core.Definition, signingIdentity *core.SignerRef, tag string, waitConfirm bool) (msg *core.Message, err error)
	CreateDatatype(ctx context.Context, datatype *core.Datatype, waitConfirm bool) (msg *core.Message, err error)
	CreateIdentityClaim(ctx context.Context, def *core.IdentityClaim, signingIdentity *core.SignerRef, tag string, waitConfirm bool) (msg *core.Message, err error)
	CreateTokenPool(ctx context.Context, pool *core.TokenPoolAnnouncement, waitConfirm bool) (msg *core.Message, err error)
}

type definitionSender struct {
	ctx        context.Context
	namespace  string
	multiparty bool
	broadcast  broadcast.Manager // optional
	identity   identity.Manager
	data       data.Manager
	handler    DefinitionHandler
}

func NewDefinitionSender(ctx context.Context, ns string, multiparty bool, bm broadcast.Manager, im identity.Manager, dm data.Manager) (Sender, error) {
	if im == nil || dm == nil {
		return nil, i18n.NewError(ctx, coremsgs.MsgInitializationNilDepError)
	}
	return &definitionSender{
		ctx:        ctx,
		namespace:  ns,
		multiparty: multiparty,
		broadcast:  bm,
		identity:   im,
		data:       dm,
	}, nil
}

func (bm *definitionSender) Name() string {
	return "DefinitionSender"
}

func (bm *definitionSender) Init(handler DefinitionHandler) {
	bm.handler = handler
}

func (bm *definitionSender) CreateDefinition(ctx context.Context, def core.Definition, tag string, waitConfirm bool) (msg *core.Message, err error) {
	return bm.CreateDefinitionWithIdentity(ctx, def, &core.SignerRef{ /* resolve to node default */ }, tag, waitConfirm)
}

func (bm *definitionSender) CreateDefinitionWithIdentity(ctx context.Context, def core.Definition, signingIdentity *core.SignerRef, tag string, waitConfirm bool) (msg *core.Message, err error) {

	if bm.multiparty {
		err = bm.identity.ResolveInputSigningIdentity(ctx, signingIdentity)
		if err != nil {
			return nil, err
		}
	}

	return bm.createDefinitionCommon(ctx, def, signingIdentity, tag, waitConfirm)
}

// CreateIdentityClaim is a special form of CreateDefinition where the signing identity does not need to have been pre-registered
// The blockchain "key" will be normalized, but the "author" will pass through unchecked
func (bm *definitionSender) CreateIdentityClaim(ctx context.Context, def *core.IdentityClaim, signingIdentity *core.SignerRef, tag string, waitConfirm bool) (msg *core.Message, err error) {

	if signingIdentity != nil {
		signingIdentity.Key, err = bm.identity.NormalizeSigningKey(ctx, signingIdentity.Key, identity.KeyNormalizationBlockchainPlugin)
		if err != nil {
			return nil, err
		}
	}

	return bm.createDefinitionCommon(ctx, def, signingIdentity, tag, waitConfirm)
}

func (bm *definitionSender) createDefinitionCommon(ctx context.Context, def core.Definition, signingIdentity *core.SignerRef, tag string, waitConfirm bool) (*core.Message, error) {

	b, err := json.Marshal(&def)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, coremsgs.MsgSerializationFailed)
	}
	dataValue := fftypes.JSONAnyPtrBytes(b)
	message := &core.MessageInOut{
		Message: core.Message{
			Header: core.MessageHeader{
				Namespace: bm.namespace,
				Type:      core.MessageTypeDefinition,
				Topics:    core.FFStringArray{def.Topic()},
				Tag:       tag,
				TxType:    core.TransactionTypeBatchPin,
			},
		},
	}
	if signingIdentity != nil {
		message.Header.SignerRef = *signingIdentity
	}

	if bm.multiparty {
		message.InlineData = core.InlineData{
			&core.DataRefOrValue{Value: dataValue},
		}
		sender := bm.broadcast.NewBroadcast(message)
		if waitConfirm {
			err = sender.SendAndWait(ctx)
		} else {
			err = sender.Send(ctx)
		}
	} else {
		data := &core.Data{Value: dataValue}
		err = bm.createDefinitionLocal(ctx, &message.Message, data)
	}

	return &message.Message, err
}

func (bm *definitionSender) createDefinitionLocal(ctx context.Context, msg *core.Message, data *core.Data) (err error) {
	state := &core.BatchState{
		PendingConfirms: make(map[fftypes.UUID]*core.Message),
	}
	if err = bm.handler.HandleDefinition(ctx, state, msg, data); err == nil {
		if err = state.RunPreFinalize(ctx); err == nil {
			err = state.RunFinalize(ctx)
		}
	}
	return err
}
