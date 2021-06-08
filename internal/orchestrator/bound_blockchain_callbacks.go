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

package orchestrator

import (
	"github.com/kaleido-io/firefly/internal/events"
	"github.com/kaleido-io/firefly/pkg/blockchain"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

type boundBlockchainCallbacks struct {
	bi blockchain.Plugin
	ei events.EventManager
}

func (bbc *boundBlockchainCallbacks) TransactionUpdate(txTrackingID string, txState blockchain.TransactionStatus, protocolTxID, errorMessage string, additionalInfo fftypes.JSONObject) error {
	return bbc.ei.TransactionUpdate(bbc.bi, txTrackingID, txState, protocolTxID, errorMessage, additionalInfo)
}

func (bbc *boundBlockchainCallbacks) BatchPinComplete(batch *blockchain.BatchPin, signingIdentity string, protocolTxID string, additionalInfo fftypes.JSONObject) error {
	return bbc.ei.BatchPinComplete(bbc.bi, batch, signingIdentity, protocolTxID, additionalInfo)
}
