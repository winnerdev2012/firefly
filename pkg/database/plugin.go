// Copyright © 2021 Kaleido, Inc.
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

	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

var (
	HashMismatch = i18n.NewError(context.Background(), i18n.MsgHashMismatch)
	IDMismatch   = i18n.NewError(context.Background(), i18n.MsgIDMismatch)
)

// Plugin is the interface implemented by each plugin
type Plugin interface {
	PeristenceInterface // Split out to aid pluggability the next level down (SQL provider etc.)

	// InitConfigPrefix initializes the set of configuration options that are valid, with defaults. Called on all plugins.
	InitConfigPrefix(prefix config.ConfigPrefix)

	// Init initializes the plugin, with configuration
	// Returns the supported featureset of the interface
	Init(ctx context.Context, prefix config.ConfigPrefix, callbacks Callbacks) error

	// Capabilities returns capabilities - not called until after Init
	Capabilities() *Capabilities
}

// The database mechanism of Firefly is designed to provide the balance between being able
// to query the data a member of the network has transferred/received via Firefly efficiently,
// while not trying to become the core database of the application (where full deeply nested
// rich query is needed).
//
// This means that we treat business data as opaque within the stroage, only verifying it against
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
type PeristenceInterface interface {
	fftypes.Named

	// RunAsGroup instructs the database plugin that all database operations performed within the context
	// function can be grouped into a single transaction (if supported).
	// Note, the caller is responsible for passing the context back to all database operations performed within the supplied function.
	RunAsGroup(ctx context.Context, fn func(ctx context.Context) error) error

	// Upsert a namespace
	// Throws IDMismatch error if updating and ids don't match
	UpsertNamespace(ctx context.Context, data *fftypes.Namespace, allowExisting bool) (err error)

	// Update namespace
	UpdateNamespace(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get an namespace by name
	GetNamespace(ctx context.Context, name string) (offset *fftypes.Namespace, err error)

	// Get namespaces
	GetNamespaces(ctx context.Context, filter Filter) (offset []*fftypes.Namespace, err error)

	// Upsert a message, with all the embedded data references.
	// allowHashUpdate=false throws HashMismatch error if the updated message has a different hash
	UpsertMessage(ctx context.Context, message *fftypes.Message, allowExisting, allowHashUpdate bool) (err error)

	// Update message
	UpdateMessage(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Update messages
	UpdateMessages(ctx context.Context, filter Filter, update Update) (err error)

	// Get a message by Id
	GetMessageById(ctx context.Context, id *fftypes.UUID) (message *fftypes.Message, err error)

	// List messages, reverse sorted (newest first) by Confirmed then Created, with pagination, and simple must filters
	GetMessages(ctx context.Context, filter Filter) (message []*fftypes.Message, err error)

	// Lighter weight query to just get the reference info of messages
	GetMessageRefs(ctx context.Context, filter Filter) ([]*fftypes.MessageRef, error)

	// List messages where there is a data reference to the specified ID
	GetMessagesForData(ctx context.Context, dataId *fftypes.UUID, filter Filter) (message []*fftypes.Message, err error)

	// Check to see if all the data for this message is available
	CheckDataAvailable(ctx context.Context, msg *fftypes.Message) (bool, error)

	// Upsert a data record
	// allowHashUpdate=false throws HashMismatch error if the updated message has a different hash
	UpsertData(ctx context.Context, data *fftypes.Data, allowExisting, allowHashUpdate bool) (err error)

	// Update data
	UpdateData(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get a data record by Id
	GetDataById(ctx context.Context, id *fftypes.UUID) (message *fftypes.Data, err error)

	// Get data
	GetData(ctx context.Context, filter Filter) (message []*fftypes.Data, err error)

	// Get data references only (no data)
	GetDataRefs(ctx context.Context, filter Filter) (message fftypes.DataRefs, err error)

	// Upsert a batch
	// allowHashUpdate=false throws HashMismatch error if the updated message has a different hash
	UpsertBatch(ctx context.Context, data *fftypes.Batch, allowExisting, allowHashUpdate bool) (err error)

	// Update data
	UpdateBatch(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get a batch by Id
	GetBatchById(ctx context.Context, id *fftypes.UUID) (message *fftypes.Batch, err error)

	// Get batches
	GetBatches(ctx context.Context, filter Filter) (message []*fftypes.Batch, err error)

	// Upsert a transaction
	// allowHashUpdate=false throws HashMismatch error if the updated message has a different hash
	UpsertTransaction(ctx context.Context, data *fftypes.Transaction, allowExisting, allowHashUpdate bool) (err error)

	// Update transaction
	UpdateTransaction(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get a transaction by Id
	GetTransactionById(ctx context.Context, id *fftypes.UUID) (message *fftypes.Transaction, err error)

	// Get transactions
	GetTransactions(ctx context.Context, filter Filter) (message []*fftypes.Transaction, err error)

	// Upsert a data definitino
	UpsertDataDefinition(ctx context.Context, datadef *fftypes.DataDefinition, allowExisting bool) (err error)

	// Update data definition
	UpdateDataDefinition(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get a data definition by Id
	GetDataDefinitionById(ctx context.Context, id *fftypes.UUID) (datadef *fftypes.DataDefinition, err error)

	// Get a data definition by name
	GetDataDefinitionByName(ctx context.Context, ns, name string) (datadef *fftypes.DataDefinition, err error)

	// Get data definitions
	GetDataDefinitions(ctx context.Context, filter Filter) (datadef []*fftypes.DataDefinition, err error)

	// Upsert an offset
	UpsertOffset(ctx context.Context, data *fftypes.Offset, allowExisting bool) (err error)

	// Update offset
	UpdateOffset(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get an offset by Id
	GetOffset(ctx context.Context, t fftypes.OffsetType, ns, name string) (offset *fftypes.Offset, err error)

	// Get offsets
	GetOffsets(ctx context.Context, filter Filter) (offset []*fftypes.Offset, err error)

	// Upsert an operation
	UpsertOperation(ctx context.Context, operation *fftypes.Operation, allowExisting bool) (err error)

	// Update matching operations
	UpdateOperations(ctx context.Context, filter Filter, update Update) (err error)

	// Get an operation by Id
	GetOperationById(ctx context.Context, id *fftypes.UUID) (operation *fftypes.Operation, err error)

	// Get operation
	GetOperations(ctx context.Context, filter Filter) (operation []*fftypes.Operation, err error)

	// Upsert a subscription
	UpsertSubscription(ctx context.Context, data *fftypes.Subscription, allowExisting bool) (err error)

	// Update subscription
	// Throws IDMismatch error if updating and ids don't match
	UpdateSubscription(ctx context.Context, ns, name string, update Update) (err error)

	// Get an subscription by name
	GetSubscription(ctx context.Context, ns, name string) (offset *fftypes.Subscription, err error)

	// Get subscriptions
	GetSubscriptions(ctx context.Context, filter Filter) (offset []*fftypes.Subscription, err error)

	// Upsert an event
	UpsertEvent(ctx context.Context, data *fftypes.Event, allowExisting bool) (err error)

	// Update event
	UpdateEvent(ctx context.Context, id *fftypes.UUID, update Update) (err error)

	// Get a event by Id
	GetEventById(ctx context.Context, id *fftypes.UUID) (message *fftypes.Event, err error)

	// Get events
	GetEvents(ctx context.Context, filter Filter) (message []*fftypes.Event, err error)
}

// Events
//
// If Capabilities returns ClusterEvents=true then these should be brodcast to every instance within
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
// TODO: Clarify the relationship between Leader Election capabilities and Event capabilities
//
type Callbacks interface {
	MessageCreated(sequence int64)
	EventCreated(sequence int64)
}

// No capabilities currently defined for the database interface - all features are mandatory
type Capabilities struct {
	ClusterEvents bool
}

var NamespaceQueryFactory = &queryFields{
	"id":          &UUIDField{},
	"type":        &StringField{},
	"name":        &StringField{},
	"description": &StringField{},
	"created":     &TimeField{},
	"confirmed":   &TimeField{},
}

var MessageQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"cid":       &UUIDField{},
	"namespace": &StringField{},
	"type":      &StringField{},
	"author":    &StringField{},
	"topic":     &StringField{},
	"context":   &StringField{},
	"group":     &UUIDField{},
	"created":   &TimeField{},
	"confirmed": &TimeField{},
	"sequence":  &Int64Field{},
	"tx.type":   &StringField{},
	"tx.id":     &UUIDField{},
	"batchid":   &UUIDField{},
}

var BatchQueryFactory = &queryFields{
	"id":         &UUIDField{},
	"namespace":  &StringField{},
	"type":       &StringField{},
	"author":     &StringField{},
	"topic":      &StringField{},
	"context":    &StringField{},
	"group":      &StringField{},
	"payloadref": &StringField{},
	"created":    &TimeField{},
	"confirmed":  &TimeField{},
	"tx.type":    &StringField{},
	"tx.id":      &UUIDField{},
}

var TransactionQueryFactory = &queryFields{
	"id":         &UUIDField{},
	"namespace":  &StringField{},
	"type":       &StringField{},
	"author":     &StringField{},
	"protocolid": &StringField{},
	"status":     &StringField{},
	"message":    &UUIDField{},
	"batch":      &UUIDField{},
	"created":    &TimeField{},
	"confirmed":  &TimeField{},
	"sequence":   &Int64Field{},
}

var DataQueryFactory = &queryFields{
	"id":                 &UUIDField{},
	"namespace":          &StringField{},
	"validator":          &StringField{},
	"definition.name":    &StringField{},
	"definition.version": &StringField{},
	"hash":               &StringField{},
	"created":            &TimeField{},
}

var DataDefinitionQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"namespace": &StringField{},
	"validator": &StringField{},
	"name":      &StringField{},
	"version":   &StringField{},
	"created":   &TimeField{},
}

var OffsetQueryFactory = &queryFields{
	"namespace": &StringField{},
	"name":      &StringField{},
	"type":      &StringField{},
	"current":   &Int64Field{},
}

var OperationQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"namespace": &StringField{},
	"message":   &UUIDField{},
	"data":      &UUIDField{},
	"type":      &StringField{},
	"recipient": &StringField{},
	"status":    &StringField{},
	"error":     &StringField{},
	"plugin":    &StringField{},
	"backendid": &StringField{},
	"created":   &TimeField{},
	"updated":   &TimeField{},
}

var SubscriptionQueryFactory = &queryFields{
	"id":             &UUIDField{},
	"namespace":      &StringField{},
	"name":           &StringField{},
	"transport":      &StringField{},
	"events":         &StringField{},
	"filter.topic":   &StringField{},
	"filter.context": &StringField{},
	"filter.group":   &StringField{},
	"options":        &StringField{},
	"created":        &TimeField{},
}

var EventQueryFactory = &queryFields{
	"id":        &UUIDField{},
	"type":      &StringField{},
	"namespace": &StringField{},
	"reference": &UUIDField{},
	"sequence":  &Int64Field{},
	"created":   &TimeField{},
}
