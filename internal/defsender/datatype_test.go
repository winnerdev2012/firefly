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
	"fmt"
	"testing"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly/mocks/broadcastmocks"
	"github.com/hyperledger/firefly/mocks/datamocks"
	"github.com/hyperledger/firefly/mocks/identitymanagermocks"
	"github.com/hyperledger/firefly/mocks/sysmessagingmocks"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateDatatypeGatewayMode(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()
	_, err := ds.CreateDatatype(context.Background(), &core.Datatype{}, false)
	assert.Regexp(t, "FF10413", err)
}

func TestCreateDatatypeBadType(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()
	ds.multiparty = true
	_, err := ds.CreateDatatype(context.Background(), &core.Datatype{
		Validator: core.ValidatorType("wrong"),
	}, false)
	assert.Regexp(t, "FF00111.*validator", err)
}

func TestBroadcastDatatypeBadValue(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()
	ds.multiparty = true
	mdm := ds.data.(*datamocks.Manager)
	mdm.On("CheckDatatype", mock.Anything, mock.Anything).Return(nil)
	mim := ds.identity.(*identitymanagermocks.Manager)
	mim.On("ResolveInputSigningIdentity", mock.Anything, mock.Anything).Return(nil)
	_, err := ds.CreateDatatype(context.Background(), &core.Datatype{
		Namespace: "ns1",
		Name:      "ent1",
		Version:   "0.0.1",
		Value:     fftypes.JSONAnyPtr(`!unparsable`),
	}, false)
	assert.Regexp(t, "FF10137.*value", err)

	mdm.AssertExpectations(t)
	mim.AssertExpectations(t)
}

func TestCreateDatatypeInvalid(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()
	ds.multiparty = true
	mdm := ds.data.(*datamocks.Manager)
	mim := ds.identity.(*identitymanagermocks.Manager)

	mim.On("ResolveInputIdentity", mock.Anything, mock.Anything).Return(nil)
	mdm.On("CheckDatatype", mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))

	_, err := ds.CreateDatatype(context.Background(), &core.Datatype{
		Namespace: "ns1",
		Name:      "ent1",
		Version:   "0.0.1",
		Value:     fftypes.JSONAnyPtr(`{"some": "data"}`),
	}, false)
	assert.EqualError(t, err, "pop")

	mdm.AssertExpectations(t)
}

func TestBroadcastOk(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()
	ds.multiparty = true
	mdm := ds.data.(*datamocks.Manager)
	mim := ds.identity.(*identitymanagermocks.Manager)
	mbm := ds.broadcast.(*broadcastmocks.Manager)
	mms := &sysmessagingmocks.MessageSender{}

	mim.On("ResolveInputSigningIdentity", mock.Anything, mock.Anything).Return(nil)
	mdm.On("CheckDatatype", mock.Anything, mock.Anything).Return(nil)
	mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("Send", context.Background()).Return(nil)

	_, err := ds.CreateDatatype(context.Background(), &core.Datatype{
		Namespace: "ns1",
		Name:      "ent1",
		Version:   "0.0.1",
		Value:     fftypes.JSONAnyPtr(`{"some": "data"}`),
	}, false)
	assert.NoError(t, err)

	mdm.AssertExpectations(t)
	mim.AssertExpectations(t)
	mbm.AssertExpectations(t)
	mms.AssertExpectations(t)
}
