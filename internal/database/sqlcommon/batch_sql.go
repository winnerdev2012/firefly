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

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/kaleido-io/firefly/internal/database"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
)

var (
	batchColumns = []string{
		"id",
		"btype",
		"namespace",
		"author",
		"created",
		"hash",
		"payload",
		"payload_ref",
		"confirmed",
		"tx_type",
		"tx_id",
	}
	batchFilterTypeMap = map[string]string{
		"type":       "btype",
		"payloadref": "payload_ref",
		"tx.type":    "tx_type",
		"tx.id":      "tx_id",
	}
)

func (s *SQLCommon) UpsertBatch(ctx context.Context, batch *fftypes.Batch, allowHashUpdate bool) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	// Do a select within the transaction to detemine if the UUID already exists
	batchRows, err := s.queryTx(ctx, tx,
		sq.Select("hash").
			From("batches").
			Where(sq.Eq{"id": batch.ID}),
	)
	if err != nil {
		return err
	}

	if batchRows.Next() {
		if !allowHashUpdate {
			var hash *fftypes.Bytes32
			_ = batchRows.Scan(&hash)
			if !fftypes.SafeHashCompare(hash, batch.Hash) {
				batchRows.Close()
				log.L(ctx).Errorf("Existing=%s New=%s", hash, batch.Hash)
				return database.HashMismatch
			}
		}
		batchRows.Close()

		// Update the batch
		if _, err = s.updateTx(ctx, tx,
			sq.Update("batches").
				Set("btype", string(batch.Type)).
				Set("namespace", batch.Namespace).
				Set("author", batch.Author).
				Set("created", batch.Created).
				Set("hash", batch.Hash).
				Set("payload", batch.Payload).
				Set("payload_ref", batch.PayloadRef).
				Set("confirmed", batch.Confirmed).
				Set("tx_type", batch.Payload.TX.Type).
				Set("tx_id", batch.Payload.TX.ID).
				Where(sq.Eq{"id": batch.ID}),
		); err != nil {
			return err
		}
	} else {
		batchRows.Close()

		if _, err = s.insertTx(ctx, tx,
			sq.Insert("batches").
				Columns(batchColumns...).
				Values(
					batch.ID,
					string(batch.Type),
					batch.Namespace,
					batch.Author,
					batch.Created,
					batch.Hash,
					batch.Payload,
					batch.PayloadRef,
					batch.Confirmed,
					batch.Payload.TX.Type,
					batch.Payload.TX.ID,
				),
		); err != nil {
			return err
		}
	}

	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) batchResult(ctx context.Context, row *sql.Rows) (*fftypes.Batch, error) {
	var batch fftypes.Batch
	err := row.Scan(
		&batch.ID,
		&batch.Type,
		&batch.Namespace,
		&batch.Author,
		&batch.Created,
		&batch.Hash,
		&batch.Payload,
		&batch.PayloadRef,
		&batch.Confirmed,
		&batch.Payload.TX.Type,
		&batch.Payload.TX.ID,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBReadErr, "batches")
	}
	return &batch, nil
}

func (s *SQLCommon) GetBatchById(ctx context.Context, ns string, id *uuid.UUID) (message *fftypes.Batch, err error) {

	rows, err := s.query(ctx,
		sq.Select(batchColumns...).
			From("batches").
			Where(sq.Eq{"id": id}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Batch '%s' not found", id)
		return nil, nil
	}

	batch, err := s.batchResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return batch, nil
}

func (s *SQLCommon) GetBatches(ctx context.Context, filter database.Filter) (message []*fftypes.Batch, err error) {

	query, err := s.filterSelect(ctx, sq.Select(batchColumns...).From("batches"), filter, batchFilterTypeMap)
	if err != nil {
		return nil, err
	}

	rows, err := s.query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := []*fftypes.Batch{}
	for rows.Next() {
		batch, err := s.batchResult(ctx, rows)
		if err != nil {
			return nil, err
		}
		batches = append(batches, batch)
	}

	return batches, err

}

func (s *SQLCommon) UpdateBatch(ctx context.Context, id *uuid.UUID, update database.Update) (err error) {

	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	query, err := s.buildUpdate(ctx, sq.Update("batches"), update, batchFilterTypeMap)
	if err != nil {
		return err
	}
	query = query.Where(sq.Eq{"id": id})

	_, err = s.updateTx(ctx, tx, query)
	if err != nil {
		return err
	}

	return s.commitTx(ctx, tx, autoCommit)
}
