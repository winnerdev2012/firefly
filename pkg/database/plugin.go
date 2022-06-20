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

package database

import (
	"context"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
)

var (
	// HashMismatch sentinel error
	HashMismatch = i18n.NewError(context.Background(), coremsgs.MsgHashMismatch)
	// IDMismatch sentinel error
	IDMismatch = i18n.NewError(context.Background(), coremsgs.MsgIDMismatch)
	// DeleteRecordNotFound sentinel error
	DeleteRecordNotFound = i18n.NewError(context.Background(), coremsgs.Msg404NotFound)
)

type UpsertOptimization int

const (
	UpsertOptimizationSkip UpsertOptimization = iota
	UpsertOptimizationNew
	UpsertOptimizationExisting
)

// Plugin is the interface implemented by each plugin
type Plugin interface {
	PersistenceInterface // Split out to aid pluggability the next level down (SQL provider etc.)

	// InitConfig initializes the set of configuration options that are valid, with defaults. Called on all plugins.
	InitConfig(config config.Section)

	// Init initializes the plugin, with configuration
	Init(ctx context.Context, config config.Section) error

	// RegisterListener registers a listener to receive callbacks
	RegisterListener(listener Callbacks)

	// Capabilities returns capabilities - not called until after Init
	Capabilities() *Capabilities
}

type iNamespaceCollection interface {
	// UpsertNamespace - Upsert a namespace
	// Throws IDMismatch error if updating and ids don't match
	UpsertNamespace(ctx context.Context, data *core.Namespace, allowExisting bool) (err error)

	// DeleteNamespace - Delete namespace
	DeleteNamespace(ctx context.Context, id *fftypes.UUID) (err error)

	// GetNamespace - Get an namespace by name
	GetNamespace(ctx context.Context, name string) (namespace *core.Namespace, err error)

	// GetNamespaceByID - Get a namespace by ID
	GetNamespaceByID(ctx context.Context, id *fftypes.UUID) (namespace *core.Namespace, err error)
}

type iMessageCollection interface {
	// UpsertMessage - Upsert a message, with all the embedded data references.
	//                 The database layer must ensure that if a record already exists, the hash of that existing record
	//                 must match the hash of the record that is being inserted.
	UpsertMessage(ctx context.Context, message *core.Message, optimization UpsertOptimization, hooks ...PostCompletionHook) (err error)

	// InsertMessages performs a batch insert of messages assured to be new records - fails if they already exist, so caller can fall back to upsert individually
	InsertMessages(ctx context.Context, messages []*core.Message, hooks ...PostCompletionHook) (err error)

	// UpdateMessage - Update message
	UpdateMessage(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// ReplaceMessage updates the message, and assigns it a new sequence number at the front of the list.
	// A new event is raised for the message, with the new sequence number - as if it was brand new.
	ReplaceMessage(ctx context.Context, message *core.Message) (err error)

	// UpdateMessages - Update messages
	UpdateMessages(ctx context.Context, filter Filter, update Update) (err error)

	// GetMessageByID - Get a message by ID
	GetMessageByID(ctx context.Context, id *fftypes.UUID) (message *core.Message, err error)

	// GetMessages - List messages, reverse sorted (newest first) by Confirmed then Created, with pagination, and simple must filters
	GetMessages(ctx context.Context, filter Filter) (message []*core.Message, res *FilterResult, err error)

	// GetMessageIDs - Retrieves messages, but only querying the messages ID (no other fields)
	GetMessageIDs(ctx context.Context, filter Filter) (ids []*core.IDAndSequence, err error)

	// GetMessagesForData - List messages where there is a data reference to the specified ID
	GetMessagesForData(ctx context.Context, dataID *fftypes.UUID, filter Filter) (message []*core.Message, res *FilterResult, err error)

	// GetBatchIDsForMessages - an optimized query to retrieve any non-null batch IDs for a list of message IDs
	GetBatchIDsForMessages(ctx context.Context, msgIDs []*fftypes.UUID) (batchIDs []*fftypes.UUID, err error)

	// GetBatchIDsForDataAttachments - an optimized query to retrieve any non-null batch IDs for a list of data IDs that might be attached to messages in batches
	GetBatchIDsForDataAttachments(ctx context.Context, dataIDs []*fftypes.UUID) (batchIDs []*fftypes.UUID, err error)
}

type iDataCollection interface {
	// UpsertData - Upsert a data record. A hint can be supplied to whether the data already exists.
	//              The database layer must ensure that if a record already exists, the hash of that existing record
	//              must match the hash of the record that is being inserted.
	UpsertData(ctx context.Context, data *core.Data, optimization UpsertOptimization) (err error)

	// InsertDataArray performs a batch insert of data assured to be new records - fails if they already exist, so caller can fall back to upsert individually
	InsertDataArray(ctx context.Context, data core.DataArray) (err error)

	// UpdateData - Update data
	UpdateData(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// GetDataByID - Get a data record by ID
	GetDataByID(ctx context.Context, id *fftypes.UUID, withValue bool) (message *core.Data, err error)

	// GetData - Get data
	GetData(ctx context.Context, filter Filter) (message core.DataArray, res *FilterResult, err error)

	// GetDataRefs - Get data references only (no data)
	GetDataRefs(ctx context.Context, filter Filter) (message core.DataRefs, res *FilterResult, err error)
}

type iBatchCollection interface {
	// UpsertBatch - Upsert a batch - the hash cannot change
	UpsertBatch(ctx context.Context, data *core.BatchPersisted) (err error)

	// UpdateBatch - Update data
	UpdateBatch(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// GetBatchByID - Get a batch by ID
	GetBatchByID(ctx context.Context, id *fftypes.UUID) (message *core.BatchPersisted, err error)

	// GetBatches - Get batches
	GetBatches(ctx context.Context, filter Filter) (message []*core.BatchPersisted, res *FilterResult, err error)
}

type iTransactionCollection interface {
	// InsertTransaction - Insert a new transaction
	InsertTransaction(ctx context.Context, data *core.Transaction) (err error)

	// UpdateTransaction - Update transaction
	UpdateTransaction(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// GetTransactionByID - Get a transaction by ID
	GetTransactionByID(ctx context.Context, id *fftypes.UUID) (message *core.Transaction, err error)

	// GetTransactions - Get transactions
	GetTransactions(ctx context.Context, filter Filter) (message []*core.Transaction, res *FilterResult, err error)
}

type iDatatypeCollection interface {
	// UpsertDatatype - Upsert a data definition
	UpsertDatatype(ctx context.Context, datadef *core.Datatype, allowExisting bool) (err error)

	// UpdateDatatype - Update data definition
	UpdateDatatype(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// GetDatatypeByID - Get a data definition by ID
	GetDatatypeByID(ctx context.Context, id *fftypes.UUID) (datadef *core.Datatype, err error)

	// GetDatatypeByName - Get a data definition by name
	GetDatatypeByName(ctx context.Context, ns, name, version string) (datadef *core.Datatype, err error)

	// GetDatatypes - Get data definitions
	GetDatatypes(ctx context.Context, filter Filter) (datadef []*core.Datatype, res *FilterResult, err error)
}

type iOffsetCollection interface {
	// UpsertOffset - Upsert an offset
	UpsertOffset(ctx context.Context, data *core.Offset, allowExisting bool) (err error)

	// UpdateOffset - Update offset
	UpdateOffset(ctx context.Context, rowID int64, update Update) (err error)

	// GetOffset - Get an offset by name
	GetOffset(ctx context.Context, t core.OffsetType, name string) (offset *core.Offset, err error)

	// GetOffsets - Get offsets
	GetOffsets(ctx context.Context, filter Filter) (offset []*core.Offset, res *FilterResult, err error)

	// DeleteOffset - Delete an offset by name
	DeleteOffset(ctx context.Context, t core.OffsetType, name string) (err error)
}

type iPinCollection interface {
	// InsertPins - Inserts a list of pins - fails if they already exist, so caller can fall back to upsert individually
	InsertPins(ctx context.Context, pins []*core.Pin) (err error)

	// UpsertPin - Will insert a pin at the end of the sequence, unless the batch+hash+index sequence already exists
	UpsertPin(ctx context.Context, parked *core.Pin) (err error)

	// GetPins - Get pins
	GetPins(ctx context.Context, filter Filter) (offset []*core.Pin, res *FilterResult, err error)

	// UpdatePins - Updates pins
	UpdatePins(ctx context.Context, filter Filter, update Update) (err error)
}

type iOperationCollection interface {
	// InsertOperation - Insert an operation
	InsertOperation(ctx context.Context, operation *core.Operation, hooks ...PostCompletionHook) (err error)

	// ResolveOperation - Resolve operation upon completion
	ResolveOperation(ctx context.Context, ns string, id *fftypes.UUID, status core.OpStatus, errorMsg *string, output fftypes.JSONObject) (err error)

	// UpdateOperation - Update an operation
	UpdateOperation(ctx context.Context, ns string, id *fftypes.UUID, update Update) (err error)

	// GetOperationByID - Get an operation by ID
	GetOperationByID(ctx context.Context, id *fftypes.UUID) (operation *core.Operation, err error)

	// GetOperations - Get operation
	GetOperations(ctx context.Context, filter Filter) (operation []*core.Operation, res *FilterResult, err error)
}

type iSubscriptionCollection interface {
	// UpsertSubscription - Upsert a subscription
	UpsertSubscription(ctx context.Context, data *core.Subscription, allowExisting bool) (err error)

	// UpdateSubscription - Update subscription
	// Throws IDMismatch error if updating and ids don't match
	UpdateSubscription(ctx context.Context, ns, name string, update Update) (err error)

	// GetSubscriptionByName - Get an subscription by name
	GetSubscriptionByName(ctx context.Context, ns, name string) (offset *core.Subscription, err error)

	// GetSubscriptionByID - Get an subscription by id
	GetSubscriptionByID(ctx context.Context, id *fftypes.UUID) (offset *core.Subscription, err error)

	// GetSubscriptions - Get subscriptions
	GetSubscriptions(ctx context.Context, filter Filter) (offset []*core.Subscription, res *FilterResult, err error)

	// DeleteSubscriptionByID - Delete a subscription
	DeleteSubscriptionByID(ctx context.Context, id *fftypes.UUID) (err error)
}

type iEventCollection interface {
	// InsertEvent - Insert an event. The order of the sequences added to the database, must match the order that
	//               the rows/objects appear available to the event dispatcher. For a concurrency enabled database
	//               with multi-operation transactions (like PSQL or other enterprise SQL based DB) we need
	//               to hold an exclusive table lock.
	InsertEvent(ctx context.Context, data *core.Event) (err error)

	// UpdateEvent - Update event
	UpdateEvent(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// GetEventByID - Get a event by ID
	GetEventByID(ctx context.Context, id *fftypes.UUID) (message *core.Event, err error)

	// GetEvents - Get events
	GetEvents(ctx context.Context, filter Filter) (message []*core.Event, res *FilterResult, err error)
}

type iIdentitiesCollection interface {
	// UpsertIdentity - Upsert an identity
	UpsertIdentity(ctx context.Context, data *core.Identity, optimization UpsertOptimization) (err error)

	// GetIdentityByDID - Get a identity by DID
	GetIdentityByDID(ctx context.Context, namespace, did string) (org *core.Identity, err error)

	// GetIdentityByName - Get a identity by name
	GetIdentityByName(ctx context.Context, iType core.IdentityType, namespace, name string) (org *core.Identity, err error)

	// GetIdentityByID - Get a identity by ID
	GetIdentityByID(ctx context.Context, namespace string, id *fftypes.UUID) (org *core.Identity, err error)

	// GetIdentities - Get identities
	GetIdentities(ctx context.Context, namespace string, filter Filter) (org []*core.Identity, res *FilterResult, err error)
}

type iVerifiersCollection interface {
	// UpsertVerifier - Upsert an verifier
	UpsertVerifier(ctx context.Context, data *core.Verifier, optimization UpsertOptimization) (err error)

	// GetVerifierByValue - Get a verifier by name
	GetVerifierByValue(ctx context.Context, vType core.VerifierType, namespace, value string) (org *core.Verifier, err error)

	// GetVerifierByHash - Get a verifier by its hash
	GetVerifierByHash(ctx context.Context, namespace string, hash *fftypes.Bytes32) (org *core.Verifier, err error)

	// GetVerifiers - Get verifiers
	GetVerifiers(ctx context.Context, namespace string, filter Filter) (org []*core.Verifier, res *FilterResult, err error)
}

type iGroupCollection interface {
	// UpsertGroup - Upsert a group, with a hint to whether to optmize for existing or new
	UpsertGroup(ctx context.Context, data *core.Group, optimization UpsertOptimization) (err error)

	// GetGroupByHash - Get a group by ID
	GetGroupByHash(ctx context.Context, namespace string, hash *fftypes.Bytes32) (node *core.Group, err error)

	// GetGroups - Get groups
	GetGroups(ctx context.Context, namespace string, filter Filter) (node []*core.Group, res *FilterResult, err error)
}

type iNonceCollection interface {
	// InsertNonce - Inserts a new nonce. Caller (batch processor) is responsible for ensuring it is the only active thread charge of assigning nonces to this context
	InsertNonce(ctx context.Context, nonce *core.Nonce) (err error)

	// UpdateNonce - Updates an existing nonce. Caller (batch processor) is responsible for ensuring it is the only active thread charge of assigning nonces to this context
	UpdateNonce(ctx context.Context, nonce *core.Nonce) (err error)

	// GetNonce - Get a context by hash
	GetNonce(ctx context.Context, hash *fftypes.Bytes32) (message *core.Nonce, err error)

	// GetNonces - Get contexts
	GetNonces(ctx context.Context, filter Filter) (node []*core.Nonce, res *FilterResult, err error)

	// DeleteNonce - Delete context by hash
	DeleteNonce(ctx context.Context, hash *fftypes.Bytes32) (err error)
}

type iNextPinCollection interface {
	// InsertNextPin - insert a nextpin
	InsertNextPin(ctx context.Context, nextpin *core.NextPin) (err error)

	// GetNextPinByContextAndIdentity - lookup nextpin by context+identity
	GetNextPinByContextAndIdentity(ctx context.Context, context *fftypes.Bytes32, identity string) (message *core.NextPin, err error)

	// GetNextPinByHash - lookup nextpin by its hash
	GetNextPinByHash(ctx context.Context, hash *fftypes.Bytes32) (message *core.NextPin, err error)

	// GetNextPins - get nextpins
	GetNextPins(ctx context.Context, filter Filter) (message []*core.NextPin, res *FilterResult, err error)

	// UpdateNextPin - update a next hash using its local database ID
	UpdateNextPin(ctx context.Context, sequence int64, update Update) (err error)

	// DeleteNextPin - delete a next hash, using its local database ID
	DeleteNextPin(ctx context.Context, sequence int64) (err error)
}

type iBlobCollection interface {
	// InsertBlob - insert a blob
	InsertBlob(ctx context.Context, blob *core.Blob) (err error)

	// InsertBlobs performs a batch insert of blobs assured to be new records - fails if they already exist, so caller can fall back to upsert individually
	InsertBlobs(ctx context.Context, blobs []*core.Blob) (err error)

	// GetBlobMatchingHash - lookup first blob batching a hash
	GetBlobMatchingHash(ctx context.Context, hash *fftypes.Bytes32) (message *core.Blob, err error)

	// GetBlobs - get blobs
	GetBlobs(ctx context.Context, filter Filter) (message []*core.Blob, res *FilterResult, err error)

	// DeleteBlob - delete a blob, using its local database ID
	DeleteBlob(ctx context.Context, sequence int64) (err error)
}

type iConfigRecordCollection interface {
	// UpsertConfigRecord - Upsert a config record
	// Throws IDMismatch error if updating and ids don't match
	UpsertConfigRecord(ctx context.Context, data *fftypes.ConfigRecord, allowExisting bool) (err error)

	// GetConfigRecord - Get a config record by key
	GetConfigRecord(ctx context.Context, key string) (offset *fftypes.ConfigRecord, err error)

	// GetConfigRecords - Get config records
	GetConfigRecords(ctx context.Context, filter Filter) (offset []*fftypes.ConfigRecord, res *FilterResult, err error)

	// DeleteConfigRecord - Delete config record
	DeleteConfigRecord(ctx context.Context, key string) (err error)
}

type iTokenPoolCollection interface {
	// UpsertTokenPool - Upsert a token pool
	UpsertTokenPool(ctx context.Context, pool *core.TokenPool) error

	// GetTokenPool - Get a token pool by name
	GetTokenPool(ctx context.Context, ns, name string) (*core.TokenPool, error)

	// GetTokenPoolByID - Get a token pool by pool ID
	GetTokenPoolByID(ctx context.Context, id *fftypes.UUID) (*core.TokenPool, error)

	// GetTokenPoolByID - Get a token pool by locator
	GetTokenPoolByLocator(ctx context.Context, connector, locator string) (*core.TokenPool, error)

	// GetTokenPools - Get token pools
	GetTokenPools(ctx context.Context, filter Filter) ([]*core.TokenPool, *FilterResult, error)
}

type iTokenBalanceCollection interface {
	// UpdateTokenBalances - Move some token balance from one account to another
	UpdateTokenBalances(ctx context.Context, transfer *core.TokenTransfer) error

	// GetTokenBalance - Get a token balance by pool and account identity
	GetTokenBalance(ctx context.Context, poolID *fftypes.UUID, tokenIndex, identity string) (*core.TokenBalance, error)

	// GetTokenBalances - Get token balances
	GetTokenBalances(ctx context.Context, filter Filter) ([]*core.TokenBalance, *FilterResult, error)

	// GetTokenAccounts - Get token accounts (all distinct addresses that have a balance)
	GetTokenAccounts(ctx context.Context, filter Filter) ([]*core.TokenAccount, *FilterResult, error)

	// GetTokenAccountPools - Get the list of pools referenced by a given account
	GetTokenAccountPools(ctx context.Context, key string, filter Filter) ([]*core.TokenAccountPool, *FilterResult, error)
}

type iTokenTransferCollection interface {
	// UpsertTokenTransfer - Upsert a token transfer
	UpsertTokenTransfer(ctx context.Context, transfer *core.TokenTransfer) error

	// GetTokenTransferByID - Get a token transfer by ID
	GetTokenTransferByID(ctx context.Context, localID *fftypes.UUID) (*core.TokenTransfer, error)

	// GetTokenTransferByProtocolID - Get a token transfer by protocol ID
	GetTokenTransferByProtocolID(ctx context.Context, connector, protocolID string) (*core.TokenTransfer, error)

	// GetTokenTransfers - Get token transfers
	GetTokenTransfers(ctx context.Context, filter Filter) ([]*core.TokenTransfer, *FilterResult, error)
}

type iTokenApprovalCollection interface {
	// UpsertTokenApproval - Upsert a token approval
	UpsertTokenApproval(ctx context.Context, approval *core.TokenApproval) error

	// UpdateTokenApprovals - Update multiple token approvals
	UpdateTokenApprovals(ctx context.Context, filter Filter, update Update) (err error)

	// GetTokenApprovalByID - Get a token approval by ID
	GetTokenApprovalByID(ctx context.Context, localID *fftypes.UUID) (*core.TokenApproval, error)

	// GetTokenTransferByProtocolID - Get a token approval by protocol ID
	GetTokenApprovalByProtocolID(ctx context.Context, connector, protocolID string) (*core.TokenApproval, error)

	// GetTokenApprovals - Get token approvals
	GetTokenApprovals(ctx context.Context, filter Filter) ([]*core.TokenApproval, *FilterResult, error)
}

type iFFICollection interface {
	UpsertFFI(ctx context.Context, cd *core.FFI) error
	GetFFIs(ctx context.Context, ns string, filter Filter) ([]*core.FFI, *FilterResult, error)
	GetFFIByID(ctx context.Context, id *fftypes.UUID) (*core.FFI, error)
	GetFFI(ctx context.Context, ns, name, version string) (*core.FFI, error)
}

type iFFIMethodCollection interface {
	UpsertFFIMethod(ctx context.Context, method *core.FFIMethod) error
	GetFFIMethod(ctx context.Context, ns string, interfaceID *fftypes.UUID, pathName string) (*core.FFIMethod, error)
	GetFFIMethods(ctx context.Context, filter Filter) (methods []*core.FFIMethod, res *FilterResult, err error)
}

type iFFIEventCollection interface {
	UpsertFFIEvent(ctx context.Context, method *core.FFIEvent) error
	GetFFIEvent(ctx context.Context, ns string, interfaceID *fftypes.UUID, pathName string) (*core.FFIEvent, error)
	GetFFIEventByID(ctx context.Context, id *fftypes.UUID) (*core.FFIEvent, error)
	GetFFIEvents(ctx context.Context, filter Filter) (events []*core.FFIEvent, res *FilterResult, err error)
}

type iContractAPICollection interface {
	UpsertContractAPI(ctx context.Context, cd *core.ContractAPI) error
	GetContractAPIs(ctx context.Context, ns string, filter AndFilter) ([]*core.ContractAPI, *FilterResult, error)
	GetContractAPIByID(ctx context.Context, id *fftypes.UUID) (*core.ContractAPI, error)
	GetContractAPIByName(ctx context.Context, ns, name string) (*core.ContractAPI, error)
}

type iContractListenerCollection interface {
	// InsertContractListener - upsert a subscription to an external smart contract
	InsertContractListener(ctx context.Context, sub *core.ContractListener) (err error)

	// GetContractListener - get smart contract subscription by name
	GetContractListener(ctx context.Context, ns, name string) (sub *core.ContractListener, err error)

	// GetContractListenerByID - get smart contract subscription by ID
	GetContractListenerByID(ctx context.Context, id *fftypes.UUID) (sub *core.ContractListener, err error)

	// GetContractListenerByBackendID - get smart contract subscription by backend ID
	GetContractListenerByBackendID(ctx context.Context, id string) (sub *core.ContractListener, err error)

	// GetContractListeners - get smart contract subscriptions
	GetContractListeners(ctx context.Context, filter Filter) ([]*core.ContractListener, *FilterResult, error)

	// DeleteContractListener - delete a subscription to an external smart contract
	DeleteContractListenerByID(ctx context.Context, id *fftypes.UUID) (err error)
}

type iBlockchainEventCollection interface {
	// InsertBlockchainEvent - insert an event from the blockchain
	InsertBlockchainEvent(ctx context.Context, event *core.BlockchainEvent) (err error)

	// GetBlockchainEventByID - get blockchain event by ID
	GetBlockchainEventByID(ctx context.Context, id *fftypes.UUID) (*core.BlockchainEvent, error)

	// GetBlockchainEventByID - get blockchain event by protocol ID
	GetBlockchainEventByProtocolID(ctx context.Context, ns string, listener *fftypes.UUID, protocolID string) (*core.BlockchainEvent, error)

	// GetBlockchainEvents - get blockchain events
	GetBlockchainEvents(ctx context.Context, filter Filter) ([]*core.BlockchainEvent, *FilterResult, error)
}

// PersistenceInterface are the operations that must be implemented by a database interface plugin.
type iChartCollection interface {
	// GetChartHistogram - Get charting data for a histogram
	GetChartHistogram(ctx context.Context, ns string, intervals []core.ChartHistogramInterval, collection CollectionName) ([]*core.ChartHistogram, error)
}

// PeristenceInterface are the operations that must be implemented by a database interfavce plugin.
// The database mechanism of Firefly is designed to provide the balance between being able
// to query the data a member of the network has transferred/received via Firefly efficiently,
// while not trying to become the core database of the application (where full deeply nested
// rich query is needed).
//
// This means that we treat business data as opaque within the storage, only verifying it against
// a data definition within the Firefly core runtime itself.
// The data types, indexes and relationships are designed to be simple, and map closely to the
// REST semantics of the Firefly API itself.
//
// As a result, the database interface could be implemented efficiently by most database technologies.
// Including both Relational/SQL and Document/NoSQL database technologies.
//
// As such we suggest the factors in choosing your database should be non-functional, such as:
// - Which provides you with the HA/DR capabilities you require
// - Which is most familiar within your existing devops pipeline for the application
// - Whether you can consolidate the HA/DR and server infrastructure for your app DB with the Firefly DB
//
// Each database does need an update to the core codebase, to provide a plugin that implements this
// interface.
// For SQL databases the process of adding a new database is simplified via the common SQL layer.
// For NoSQL databases, the code should be straight forward to map the collections, indexes, and operations.
//
type PersistenceInterface interface {
	core.Named

	// RunAsGroup instructs the database plugin that all database operations performed within the context
	// function can be grouped into a single transaction (if supported).
	// Requirements:
	// - Firefly must not depend on this to guarantee ACID properties (it is only a suggestion/optimization)
	// - The database implementation must support nested RunAsGroup calls (ie by reusing a transaction if one exists)
	// - The caller is responsible for passing the supplied context to all database operations within the callback function
	RunAsGroup(ctx context.Context, fn func(ctx context.Context) error) error

	iNamespaceCollection
	iMessageCollection
	iDataCollection
	iBatchCollection
	iTransactionCollection
	iDatatypeCollection
	iOffsetCollection
	iPinCollection
	iOperationCollection
	iSubscriptionCollection
	iEventCollection
	iIdentitiesCollection
	iVerifiersCollection
	iGroupCollection
	iNonceCollection
	iNextPinCollection
	iBlobCollection
	iConfigRecordCollection
	iTokenPoolCollection
	iTokenBalanceCollection
	iTokenTransferCollection
	iTokenApprovalCollection
	iFFICollection
	iFFIMethodCollection
	iFFIEventCollection
	iContractAPICollection
	iContractListenerCollection
	iBlockchainEventCollection
	iChartCollection
}

// CollectionName represents all collections
type CollectionName string

// OrderedUUIDCollectionNS collections have a strong order that includes a sequence integer
// that uniquely identifies the entry in a sequence. The sequence is LOCAL to this
// FireFly node. We try to minimize adding new collections of this type, as they have
// implementation complexity in some databases (such as NoSQL databases)
type OrderedUUIDCollectionNS CollectionName

const (
	CollectionMessages OrderedUUIDCollectionNS = "messages"
	CollectionEvents   OrderedUUIDCollectionNS = "events"
)

// OrderedCollectionNS is a collection that is ordered, and that sequence is the only key
type OrderedCollectionNS CollectionName

const (
	CollectionPins OrderedCollectionNS = "pins"
)

// UUIDCollectionNS is the most common type of collection - each entry has a UUID that
// is globally unique, and used externally by apps to address entries in the collection.
// Objects in these collections are all namespaced,.
type UUIDCollectionNS CollectionName

const (
	CollectionBatches           UUIDCollectionNS = "batches"
	CollectionBlockchainEvents  UUIDCollectionNS = "blockchainevents"
	CollectionData              UUIDCollectionNS = "data"
	CollectionDataTypes         UUIDCollectionNS = "datatypes"
	CollectionOperations        UUIDCollectionNS = "operations"
	CollectionSubscriptions     UUIDCollectionNS = "subscriptions"
	CollectionTransactions      UUIDCollectionNS = "transactions"
	CollectionTokenPools        UUIDCollectionNS = "tokenpools"
	CollectionTokenTransfers    UUIDCollectionNS = "tokentransfers"
	CollectionTokenApprovals    UUIDCollectionNS = "tokenapprovals"
	CollectionFFIs              UUIDCollectionNS = "ffi"
	CollectionFFIMethods        UUIDCollectionNS = "ffimethods"
	CollectionFFIEvents         UUIDCollectionNS = "ffievents"
	CollectionContractAPIs      UUIDCollectionNS = "contractapis"
	CollectionContractListeners UUIDCollectionNS = "contractlisteners"
	CollectionIdentities        UUIDCollectionNS = "identities"
)

// HashCollectionNS is a collection where the primary key is a hash, such that it can
// by identified by any member of the network at any time, without it first having
// been broadcast.
type HashCollectionNS CollectionName

const (
	CollectionGroups    HashCollectionNS = "groups"
	CollectionVerifiers HashCollectionNS = "verifiers"
)

// UUIDCollection is like UUIDCollectionNS, but for objects that do not reside within a namespace
type UUIDCollection CollectionName

const (
	CollectionNamespaces UUIDCollection = "namespaces"
)

// OtherCollection are odd balls, that don't fit any of the categories above.
// These collections do not support change events, and generally their
// creation is coordinated with creation of another object that does support change events.
// Mainly they are entries that require lookup by compound IDs.
type OtherCollection CollectionName

const (
	CollectionConfigrecords OtherCollection = "configrecords"
	CollectionBlobs         OtherCollection = "blobs"
	CollectionNextpins      OtherCollection = "nextpins"
	CollectionNonces        OtherCollection = "nonces"
	CollectionOffsets       OtherCollection = "offsets"
	CollectionTokenBalances OtherCollection = "tokenbalances"
)

// PostCompletionHook is a closure/function that will be called after a successful insertion.
// This includes where the insert is nested in a RunAsGroup, and the database is transactional.
// These hooks are useful when triggering code that relies on the inserted database object being available.
type PostCompletionHook func()

// Callbacks are the methods for passing data from plugin to core
//
// If Capabilities returns ClusterEvents=true then these should be broadcast to every instance within
// a cluster that is connected to the database.
//
// If Capabilities returns ClusterEvents=false then these events can be simply coupled in-process to
// update activities.
//
// The system does not rely on these events exclusively for data/transaction integrity, but if an event is
// missed/delayed it might result in slower processing.
// For example, the batch interface will initiate a batch as soon as an event is triggered, but it will use
// a subsequent database query as the source of truth of the latest set/order of data, and it will periodically
// check for new messages even if it does not receive any events.
//
// Events are emitted locally to the individual FireFly core process. However, a WebSocket interface is
// available for remote listening to these events. That allows the UI to listen to the events, as well as
// providing a building block for a cluster of FireFly servers to directly propgate events to each other.
//
type Callbacks interface {
	// OrderedUUIDCollectionNSEvent emits the sequence on insert, but it will be -1 on update
	OrderedUUIDCollectionNSEvent(resType OrderedUUIDCollectionNS, eventType core.ChangeEventType, ns string, id *fftypes.UUID, sequence int64)
	OrderedCollectionNSEvent(resType OrderedCollectionNS, eventType core.ChangeEventType, ns string, sequence int64)
	UUIDCollectionNSEvent(resType UUIDCollectionNS, eventType core.ChangeEventType, ns string, id *fftypes.UUID)
	UUIDCollectionEvent(resType UUIDCollection, eventType core.ChangeEventType, id *fftypes.UUID)
	HashCollectionNSEvent(resType HashCollectionNS, eventType core.ChangeEventType, ns string, hash *fftypes.Bytes32)
}

// Capabilities defines the capabilities a plugin can report as implementing or not
type Capabilities struct {
	Concurrency bool
}

// MessageQueryFactory filter fields for messages
var MessageQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"cid":       &UUIDField{},
	"namespace": &StringField{},
	"type":      &StringField{},
	"author":    &StringField{},
	"key":       &StringField{},
	"topics":    &FFStringArrayField{},
	"tag":       &StringField{},
	"group":     &Bytes32Field{},
	"created":   &TimeField{},
	"hash":      &Bytes32Field{},
	"pins":      &FFStringArrayField{},
	"state":     &StringField{},
	"confirmed": &TimeField{},
	"sequence":  &Int64Field{},
	"txtype":    &StringField{},
	"batch":     &UUIDField{},
}

// BatchQueryFactory filter fields for batches
var BatchQueryFactory = &queryFields{
	"id":         &UUIDField{},
	"namespace":  &StringField{},
	"type":       &StringField{},
	"author":     &StringField{},
	"key":        &StringField{},
	"group":      &Bytes32Field{},
	"hash":       &Bytes32Field{},
	"payloadref": &StringField{},
	"created":    &TimeField{},
	"confirmed":  &TimeField{},
	"tx.type":    &StringField{},
	"tx.id":      &UUIDField{},
	"node":       &UUIDField{},
}

// TransactionQueryFactory filter fields for transactions
var TransactionQueryFactory = &queryFields{
	"id":            &UUIDField{},
	"type":          &StringField{},
	"created":       &TimeField{},
	"namespace":     &StringField{},
	"blockchainids": &FFStringArrayField{},
}

// DataQueryFactory filter fields for data
var DataQueryFactory = &queryFields{
	"id":               &UUIDField{},
	"namespace":        &StringField{},
	"validator":        &StringField{},
	"datatype.name":    &StringField{},
	"datatype.version": &StringField{},
	"hash":             &Bytes32Field{},
	"blob.hash":        &Bytes32Field{},
	"blob.public":      &StringField{},
	"blob.name":        &StringField{},
	"blob.size":        &Int64Field{},
	"created":          &TimeField{},
	"value":            &JSONField{},
}

// DatatypeQueryFactory filter fields for data definitions
var DatatypeQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"message":   &UUIDField{},
	"namespace": &StringField{},
	"validator": &StringField{},
	"name":      &StringField{},
	"version":   &StringField{},
	"created":   &TimeField{},
}

// OffsetQueryFactory filter fields for data offsets
var OffsetQueryFactory = &queryFields{
	"name":    &StringField{},
	"type":    &StringField{},
	"current": &Int64Field{},
}

// OperationQueryFactory filter fields for data operations
var OperationQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"tx":        &UUIDField{},
	"type":      &StringField{},
	"namespace": &StringField{},
	"status":    &StringField{},
	"error":     &StringField{},
	"plugin":    &StringField{},
	"input":     &JSONField{},
	"output":    &JSONField{},
	"created":   &TimeField{},
	"updated":   &TimeField{},
	"retry":     &UUIDField{},
}

// SubscriptionQueryFactory filter fields for data subscriptions
var SubscriptionQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"namespace": &StringField{},
	"name":      &StringField{},
	"transport": &StringField{},
	"events":    &StringField{},
	"filters":   &JSONField{},
	"options":   &StringField{},
	"created":   &TimeField{},
}

// EventQueryFactory filter fields for data events
var EventQueryFactory = &queryFields{
	"id":         &UUIDField{},
	"type":       &StringField{},
	"namespace":  &StringField{},
	"reference":  &UUIDField{},
	"correlator": &UUIDField{},
	"tx":         &UUIDField{},
	"topic":      &StringField{},
	"sequence":   &Int64Field{},
	"created":    &TimeField{},
}

// PinQueryFactory filter fields for parked contexts
var PinQueryFactory = &queryFields{
	"namespace":  &StringField{},
	"sequence":   &Int64Field{},
	"masked":     &BoolField{},
	"hash":       &Bytes32Field{},
	"batch":      &UUIDField{},
	"index":      &Int64Field{},
	"dispatched": &BoolField{},
	"created":    &TimeField{},
}

// IdentityQueryFactory filter fields for identities
var IdentityQueryFactory = &queryFields{
	"id":                    &UUIDField{},
	"did":                   &StringField{},
	"parent":                &UUIDField{},
	"messages.claim":        &UUIDField{},
	"messages.verification": &UUIDField{},
	"messages.update":       &UUIDField{},
	"type":                  &StringField{},
	"namespace":             &StringField{},
	"name":                  &StringField{},
	"description":           &StringField{},
	"profile":               &JSONField{},
	"created":               &TimeField{},
	"updated":               &TimeField{},
}

// VerifierQueryFactory filter fields for identities
var VerifierQueryFactory = &queryFields{
	"hash":      &Bytes32Field{},
	"identity":  &UUIDField{},
	"type":      &StringField{},
	"namespace": &StringField{},
	"value":     &StringField{},
	"created":   &TimeField{},
}

// GroupQueryFactory filter fields for nodes
var GroupQueryFactory = &queryFields{
	"hash":        &Bytes32Field{},
	"message":     &UUIDField{},
	"namespace":   &StringField{},
	"description": &StringField{},
	"ledger":      &UUIDField{},
	"created":     &TimeField{},
}

// NonceQueryFactory filter fields for nodes
var NonceQueryFactory = &queryFields{
	"hash":  &StringField{},
	"nonce": &Int64Field{},
}

// NextPinQueryFactory filter fields for nodes
var NextPinQueryFactory = &queryFields{
	"context":  &Bytes32Field{},
	"identity": &StringField{},
	"hash":     &Bytes32Field{},
	"nonce":    &Int64Field{},
}

// ConfigRecordQueryFactory filter fields for config records
var ConfigRecordQueryFactory = &queryFields{
	"key":   &StringField{},
	"value": &StringField{},
}

// BlobQueryFactory filter fields for config records
var BlobQueryFactory = &queryFields{
	"hash":       &Bytes32Field{},
	"size":       &Int64Field{},
	"payloadref": &StringField{},
	"created":    &TimeField{},
}

// TokenPoolQueryFactory filter fields for token pools
var TokenPoolQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"type":      &StringField{},
	"namespace": &StringField{},
	"name":      &StringField{},
	"standard":  &StringField{},
	"locator":   &StringField{},
	"symbol":    &StringField{},
	"decimals":  &Int64Field{},
	"message":   &UUIDField{},
	"state":     &StringField{},
	"created":   &TimeField{},
	"connector": &StringField{},
	"tx.type":   &StringField{},
	"tx.id":     &UUIDField{},
}

// TokenBalanceQueryFactory filter fields for token balances
var TokenBalanceQueryFactory = &queryFields{
	"pool":       &UUIDField{},
	"tokenindex": &StringField{},
	"uri":        &StringField{},
	"connector":  &StringField{},
	"namespace":  &StringField{},
	"key":        &StringField{},
	"balance":    &Int64Field{},
	"updated":    &TimeField{},
}

// TokenAccountQueryFactory filter fields for token accounts
var TokenAccountQueryFactory = &queryFields{
	"key":       &StringField{},
	"namespace": &StringField{},
	"updated":   &TimeField{},
}

// TokenAccountPoolQueryFactory filter fields for token account pools
var TokenAccountPoolQueryFactory = &queryFields{
	"pool":      &UUIDField{},
	"namespace": &StringField{},
	"updated":   &TimeField{},
}

// TokenTransferQueryFactory filter fields for token transfers
var TokenTransferQueryFactory = &queryFields{
	"localid":         &StringField{},
	"pool":            &UUIDField{},
	"tokenindex":      &StringField{},
	"uri":             &StringField{},
	"connector":       &StringField{},
	"namespace":       &StringField{},
	"key":             &StringField{},
	"from":            &StringField{},
	"to":              &StringField{},
	"amount":          &Int64Field{},
	"protocolid":      &StringField{},
	"message":         &UUIDField{},
	"messagehash":     &Bytes32Field{},
	"created":         &TimeField{},
	"tx.type":         &StringField{},
	"tx.id":           &UUIDField{},
	"blockchainevent": &UUIDField{},
	"type":            &StringField{},
}

var TokenApprovalQueryFactory = &queryFields{
	"localid":         &StringField{},
	"pool":            &UUIDField{},
	"connector":       &StringField{},
	"namespace":       &StringField{},
	"key":             &StringField{},
	"operator":        &StringField{},
	"approved":        &BoolField{},
	"protocolid":      &StringField{},
	"subject":         &StringField{},
	"active":          &BoolField{},
	"created":         &TimeField{},
	"tx.type":         &StringField{},
	"tx.id":           &UUIDField{},
	"blockchainevent": &UUIDField{},
}

// FFIQueryFactory filter fields for contract definitions
var FFIQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"namespace": &StringField{},
	"name":      &StringField{},
	"version":   &StringField{},
}

// FFIMethodQueryFactory filter fields for contract methods
var FFIMethodQueryFactory = &queryFields{
	"id":          &UUIDField{},
	"namespace":   &StringField{},
	"name":        &StringField{},
	"pathname":    &StringField{},
	"interface":   &UUIDField{},
	"description": &StringField{},
}

// FFIEventQueryFactory filter fields for contract events
var FFIEventQueryFactory = &queryFields{
	"id":          &UUIDField{},
	"namespace":   &StringField{},
	"name":        &StringField{},
	"pathname":    &StringField{},
	"interface":   &UUIDField{},
	"description": &StringField{},
}

// ContractListenerQueryFactory filter fields for contract listeners
var ContractListenerQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"interface": &UUIDField{},
	"namespace": &StringField{},
	"location":  &JSONField{},
	"topic":     &StringField{},
	"signature": &StringField{},
	"backendid": &StringField{},
	"created":   &TimeField{},
	"updated":   &TimeField{},
	"state":     &JSONField{},
}

// BlockchainEventQueryFactory filter fields for contract events
var BlockchainEventQueryFactory = &queryFields{
	"id":              &UUIDField{},
	"source":          &StringField{},
	"namespace":       &StringField{},
	"name":            &StringField{},
	"protocolid":      &StringField{},
	"listener":        &StringField{},
	"tx.type":         &StringField{},
	"tx.id":           &UUIDField{},
	"tx.blockchainid": &StringField{},
	"timestamp":       &TimeField{},
}

// ContractAPIQueryFactory filter fields for Contract APIs
var ContractAPIQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"name":      &StringField{},
	"namespace": &StringField{},
	"interface": &UUIDField{},
}
