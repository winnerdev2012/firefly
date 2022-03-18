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

package fftypes

type TokenTransferType = FFEnum

var (
	TokenTransferTypeMint     = ffEnum("tokentransfertype", "mint")
	TokenTransferTypeBurn     = ffEnum("tokentransfertype", "burn")
	TokenTransferTypeTransfer = ffEnum("tokentransfertype", "transfer")
)

type TokenTransfer struct {
	Type            TokenTransferType `json:"type" ffenum:"tokentransfertype"`
	LocalID         *UUID             `json:"localId,omitempty"`
	Pool            *UUID             `json:"pool,omitempty"`
	TokenIndex      string            `json:"tokenIndex,omitempty"`
	URI             string            `json:"uri,omitempty"`
	Connector       string            `json:"connector,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	Key             string            `json:"key,omitempty"`
	From            string            `json:"from,omitempty"`
	To              string            `json:"to,omitempty"`
	Amount          FFBigInt          `json:"amount"`
	ProtocolID      string            `json:"protocolId,omitempty"`
	Message         *UUID             `json:"message,omitempty"`
	MessageHash     *Bytes32          `json:"messageHash,omitempty"`
	Created         *FFTime           `json:"created,omitempty"`
	TX              TransactionRef    `json:"tx"`
	BlockchainEvent *UUID             `json:"blockchainEvent,omitempty"`
}

type TokenTransferInput struct {
	TokenTransfer
	Message *MessageInOut `json:"message,omitempty"`
	Pool    string        `json:"pool,omitempty"`
}
