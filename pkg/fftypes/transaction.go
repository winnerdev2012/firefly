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

type TransactionType = FFEnum

var (
	// TransactionTypeNone indicates no transaction should be used for this message/batch
	TransactionTypeNone TransactionType = ffEnum("txtype", "none")
	// TransactionTypeBatchPin represents a pinning transaction, that verifies the originator of the data, and sequences the event deterministically between parties
	TransactionTypeBatchPin TransactionType = ffEnum("txtype", "batch_pin")
	// TransactionTypeTokenPool represents a token pool creation
	TransactionTypeTokenPool TransactionType = ffEnum("txtype", "token_pool")
	// TransactionTypeTokenTransfer represents a token transfer
	TransactionTypeTokenTransfer TransactionType = ffEnum("txtype", "token_transfer")
)

// TransactionRef refers to a transaction, in other types
type TransactionRef struct {
	Type TransactionType `json:"type"`
	ID   *UUID           `json:"id,omitempty"`
}

// Transaction represents (blockchain) transactions that were submitted by this
// node, with the correlation information to look them up on the underlying
// ledger technology
type Transaction struct {
	ID        *UUID           `json:"id,omitempty"`
	Namespace string          `json:"namespace,omitempty"`
	Type      TransactionType `json:"type" ffenum:"txtype"`
	Created   *FFTime         `json:"created"`
	Status    OpStatus        `json:"status"`
}
