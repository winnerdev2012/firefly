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

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
)

var (
	namespaceColumns = []string{
		"name",
		"description",
		"created",
		"firefly_contracts",
	}
)

const namespacesTable = "namespaces"

func (s *SQLCommon) UpsertNamespace(ctx context.Context, namespace *core.Namespace, allowExisting bool) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	existing := false
	if allowExisting {
		// Do a select within the transaction to determine if the UUID already exists
		namespaceRows, _, err := s.queryTx(ctx, namespacesTable, tx,
			sq.Select("seq").
				From(namespacesTable).
				Where(sq.Eq{"name": namespace.Name}),
		)
		if err != nil {
			return err
		}
		existing = namespaceRows.Next()
		namespaceRows.Close()
	}

	if existing {
		// Update the namespace
		if _, err = s.updateTx(ctx, namespacesTable, tx,
			sq.Update(namespacesTable).
				// Note we do not update ID
				Set("name", namespace.Name).
				Set("description", namespace.Description).
				Set("created", namespace.Created).
				Set("firefly_contracts", namespace.Contracts).
				Where(sq.Eq{"name": namespace.Name}),
			nil,
		); err != nil {
			return err
		}
	} else {
		if _, err = s.insertTx(ctx, namespacesTable, tx,
			sq.Insert(namespacesTable).
				Columns(namespaceColumns...).
				Values(
					namespace.Name,
					namespace.Description,
					namespace.Created,
					namespace.Contracts,
				),
			nil,
		); err != nil {
			return err
		}
	}

	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) namespaceResult(ctx context.Context, row *sql.Rows) (*core.Namespace, error) {
	namespace := core.Namespace{}
	err := row.Scan(
		&namespace.Name,
		&namespace.Description,
		&namespace.Created,
		&namespace.Contracts,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, coremsgs.MsgDBReadErr, namespacesTable)
	}
	return &namespace, nil
}

func (s *SQLCommon) getNamespaceEq(ctx context.Context, eq sq.Eq, textName string) (message *core.Namespace, err error) {
	rows, _, err := s.query(ctx, namespacesTable,
		sq.Select(namespaceColumns...).
			From(namespacesTable).
			Where(eq),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Namespace '%s' not found", textName)
		return nil, nil
	}

	namespace, err := s.namespaceResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (s *SQLCommon) GetNamespace(ctx context.Context, name string) (message *core.Namespace, err error) {
	return s.getNamespaceEq(ctx, sq.Eq{"name": name}, name)
}
