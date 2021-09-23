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
	"fmt"
	"testing"

	"github.com/hyperledger/firefly/mocks/blockchainmocks"
	"github.com/hyperledger/firefly/mocks/databasemocks"
	"github.com/hyperledger/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOperationUpdateSuccess(t *testing.T) {
	em, cancel := newTestEventManager(t)
	defer cancel()
	mdi := em.database.(*databasemocks.Plugin)
	mbi := &blockchainmocks.Plugin{}

	opID := fftypes.NewUUID()
	mdi.On("GetOperationByID", em.ctx, uuidMatches(opID)).Return(&fftypes.Operation{ID: opID}, nil)
	mdi.On("UpdateOperation", em.ctx, uuidMatches(opID), mock.Anything).Return(nil)

	info := fftypes.JSONObject{"some": "info"}
	err := em.OperationUpdate(mbi, opID, fftypes.OpStatusFailed, "some error", info)
	assert.NoError(t, err)

	mdi.AssertExpectations(t)
	mbi.AssertExpectations(t)
}

func TestOperationUpdateNotFound(t *testing.T) {
	em, cancel := newTestEventManager(t)
	defer cancel()
	mdi := em.database.(*databasemocks.Plugin)
	mbi := &blockchainmocks.Plugin{}

	opID := fftypes.NewUUID()
	mdi.On("GetOperationByID", em.ctx, uuidMatches(opID)).Return(nil, fmt.Errorf("pop"))

	info := fftypes.JSONObject{"some": "info"}
	err := em.OperationUpdate(mbi, opID, fftypes.OpStatusFailed, "some error", info)
	assert.NoError(t, err) // swallowed after logging

	mdi.AssertExpectations(t)
	mbi.AssertExpectations(t)
}

func TestOperationUpdateError(t *testing.T) {
	em, cancel := newTestEventManager(t)
	defer cancel()
	mdi := em.database.(*databasemocks.Plugin)
	mbi := &blockchainmocks.Plugin{}

	opID := fftypes.NewUUID()
	mdi.On("GetOperationByID", em.ctx, uuidMatches(opID)).Return(&fftypes.Operation{ID: opID}, nil)
	mdi.On("UpdateOperation", em.ctx, uuidMatches(opID), mock.Anything).Return(fmt.Errorf("pop"))

	info := fftypes.JSONObject{"some": "info"}
	err := em.OperationUpdate(mbi, opID, fftypes.OpStatusFailed, "some error", info)
	assert.EqualError(t, err, "pop")

	mdi.AssertExpectations(t)
	mbi.AssertExpectations(t)
}
