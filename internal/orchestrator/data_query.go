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

package orchestrator

import (
	"context"

	"github.com/google/uuid"
	"github.com/kaleido-io/firefly/internal/database"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
)

func (e *orchestrator) GetNamespace(ctx context.Context, ns string) (*fftypes.Namespace, error) {
	return e.database.GetNamespace(ctx, ns)
}

func (e *orchestrator) GetTransactionById(ctx context.Context, ns, id string) (*fftypes.Transaction, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgInvalidUUID)
	}
	return e.database.GetTransactionById(ctx, ns, &u)
}

func (e *orchestrator) GetMessageById(ctx context.Context, ns, id string) (*fftypes.Message, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgInvalidUUID)
	}
	return e.database.GetMessageById(ctx, ns, &u)
}

func (e *orchestrator) GetBatchById(ctx context.Context, ns, id string) (*fftypes.Batch, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgInvalidUUID)
	}
	return e.database.GetBatchById(ctx, ns, &u)
}

func (e *orchestrator) GetDataById(ctx context.Context, ns, id string) (*fftypes.Data, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgInvalidUUID)
	}
	return e.database.GetDataById(ctx, ns, &u)
}

func (e *orchestrator) GetDataDefinitionById(ctx context.Context, ns, id string) (*fftypes.DataDefinition, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgInvalidUUID)
	}
	return e.database.GetDataDefinitionById(ctx, ns, &u)
}

func (e *orchestrator) GetNamespaces(ctx context.Context, filter database.AndFilter) ([]*fftypes.Namespace, error) {
	return e.database.GetNamespaces(ctx, filter)
}

func (e *orchestrator) scopeNS(ns string, filter database.AndFilter) database.AndFilter {
	return filter.Condition(filter.Builder().Eq("namespace", ns))
}

func (e *orchestrator) GetTransactions(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Transaction, error) {
	filter = e.scopeNS(ns, filter)
	return e.database.GetTransactions(ctx, filter)
}

func (e *orchestrator) GetMessages(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Message, error) {
	filter = e.scopeNS(ns, filter)
	return e.database.GetMessages(ctx, filter)
}

func (e *orchestrator) GetMessageOperations(ctx context.Context, ns, id string, filter database.AndFilter) ([]*fftypes.Operation, error) {
	filter = e.scopeNS(ns, filter)
	filter = filter.Condition(filter.Builder().Eq("message", id))
	return e.database.GetOperations(ctx, filter)
}

func (e *orchestrator) GetBatches(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Batch, error) {
	filter = e.scopeNS(ns, filter)
	return e.database.GetBatches(ctx, filter)
}

func (e *orchestrator) GetData(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.Data, error) {
	filter = e.scopeNS(ns, filter)
	return e.database.GetData(ctx, filter)
}

func (e *orchestrator) GetDataDefinitions(ctx context.Context, ns string, filter database.AndFilter) ([]*fftypes.DataDefinition, error) {
	filter = e.scopeNS(ns, filter)
	return e.database.GetDataDefinitions(ctx, filter)
}
