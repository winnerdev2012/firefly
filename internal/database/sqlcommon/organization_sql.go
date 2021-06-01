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

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

var (
	organizationColumns = []string{
		"id",
		"message_id",
		"parent",
		"identity",
		"description",
		"profile",
		"created",
	}
	organizationFilterTypeMap = map[string]string{
		"message": "message_id",
	}
)

func (s *SQLCommon) UpsertOrganization(ctx context.Context, organization *fftypes.Organization, allowExisting bool) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	existing := false
	if allowExisting {
		// Do a select within the transaction to detemine if the UUID already exists
		organizationRows, err := s.queryTx(ctx, tx,
			sq.Select("id").
				From("orgs").
				Where(sq.Eq{"identity": organization.Identity}),
		)
		if err != nil {
			return err
		}
		existing = organizationRows.Next()

		if existing {
			var id fftypes.UUID
			_ = organizationRows.Scan(&id)
			if organization.ID != nil {
				if *organization.ID != id {
					organizationRows.Close()
					return database.IDMismatch
				}
			}
			organization.ID = &id // Update on returned object
		}
		organizationRows.Close()
	}

	if existing {
		// Update the organization
		if err = s.updateTx(ctx, tx,
			sq.Update("orgs").
				// Note we do not update ID
				Set("message_id", organization.Message).
				Set("parent", organization.Parent).
				Set("identity", organization.Identity).
				Set("description", organization.Description).
				Set("profile", organization.Profile).
				Set("created", organization.Created).
				Where(sq.Eq{"identity": organization.Identity}),
		); err != nil {
			return err
		}
	} else {
		if _, err = s.insertTx(ctx, tx,
			sq.Insert("orgs").
				Columns(organizationColumns...).
				Values(
					organization.ID,
					organization.Message,
					organization.Parent,
					organization.Identity,
					organization.Description,
					organization.Profile,
					organization.Created,
				),
		); err != nil {
			return err
		}
	}

	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) organizationResult(ctx context.Context, row *sql.Rows) (*fftypes.Organization, error) {
	organization := fftypes.Organization{}
	err := row.Scan(
		&organization.ID,
		&organization.Message,
		&organization.Parent,
		&organization.Identity,
		&organization.Description,
		&organization.Profile,
		&organization.Created,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBReadErr, "orgs")
	}
	return &organization, nil
}

func (s *SQLCommon) getOrganizationPred(ctx context.Context, desc string, pred interface{}) (message *fftypes.Organization, err error) {

	rows, err := s.query(ctx,
		sq.Select(organizationColumns...).
			From("orgs").
			Where(pred),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Organization '%s' not found", desc)
		return nil, nil
	}

	organization, err := s.organizationResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func (s *SQLCommon) GetOrganization(ctx context.Context, identity string) (message *fftypes.Organization, err error) {
	return s.getOrganizationPred(ctx, identity, sq.Eq{"identity": identity})
}

func (s *SQLCommon) GetOrganizationByID(ctx context.Context, id *fftypes.UUID) (message *fftypes.Organization, err error) {
	return s.getOrganizationPred(ctx, id.String(), sq.Eq{"id": id})
}

func (s *SQLCommon) GetOrganizations(ctx context.Context, filter database.Filter) (message []*fftypes.Organization, err error) {

	query, err := s.filterSelect(ctx, "", sq.Select(organizationColumns...).From("orgs"), filter, organizationFilterTypeMap)
	if err != nil {
		return nil, err
	}

	rows, err := s.query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	organization := []*fftypes.Organization{}
	for rows.Next() {
		d, err := s.organizationResult(ctx, rows)
		if err != nil {
			return nil, err
		}
		organization = append(organization, d)
	}

	return organization, err

}

func (s *SQLCommon) UpdateOrganization(ctx context.Context, id *fftypes.UUID, update database.Update) (err error) {

	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	query, err := s.buildUpdate(sq.Update("orgs"), update, organizationFilterTypeMap)
	if err != nil {
		return err
	}
	query = query.Where(sq.Eq{"id": id})

	err = s.updateTx(ctx, tx, query)
	if err != nil {
		return err
	}

	return s.commitTx(ctx, tx, autoCommit)
}
