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
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

var (
	subscriptionColumns = []string{
		"id",
		"namespace",
		"name",
		"transport",
		"filter_events",
		"filter_topic",
		"filter_context",
		"filter_group",
		"options",
		"created",
	}
	subscriptionFilterTypeMap = map[string]string{
		"filter.events":  "filter_events",
		"filter.topic":   "filter_topic",
		"filter.context": "filter_context",
		"filter.group":   "filter_group",
	}
)

func (s *SQLCommon) UpsertSubscription(ctx context.Context, subscription *fftypes.Subscription, allowExisting bool) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	existing := false
	if allowExisting {
		// Do a select within the transaction to detemine if the UUID already exists
		subscriptionRows, err := s.queryTx(ctx, tx,
			sq.Select("id").
				From("subscriptions").
				Where(sq.Eq{
					"namespace": subscription.Namespace,
					"name":      subscription.Name,
				}),
		)
		if err != nil {
			return err
		}

		existing = subscriptionRows.Next()
		if existing {
			var id fftypes.UUID
			_ = subscriptionRows.Scan(&id)
			if subscription.ID != nil {
				if *subscription.ID != id {
					subscriptionRows.Close()
					return database.IDMismatch
				}
			}
			subscription.ID = &id // Update on returned object
		}
		subscriptionRows.Close()
	}

	if existing {
		// Update the subscription
		if err = s.updateTx(ctx, tx,
			sq.Update("subscriptions").
				// Note we do not update ID
				Set("namespace", subscription.Namespace).
				Set("name", subscription.Name).
				Set("transport", subscription.Transport).
				Set("filter_events", subscription.Filter.Events).
				Set("filter_topic", subscription.Filter.Topic).
				Set("filter_context", subscription.Filter.Context).
				Set("filter_group", subscription.Filter.Group).
				Set("options", subscription.Options).
				Set("created", subscription.Created).
				Where(sq.Eq{
					"namespace": subscription.Namespace,
					"name":      subscription.Name,
				}),
		); err != nil {
			return err
		}
	} else {
		if subscription.ID == nil {
			subscription.ID = fftypes.NewUUID()
		}

		if _, err = s.insertTx(ctx, tx,
			sq.Insert("subscriptions").
				Columns(subscriptionColumns...).
				Values(
					subscription.ID,
					subscription.Namespace,
					subscription.Name,
					subscription.Transport,
					subscription.Filter.Events,
					subscription.Filter.Topic,
					subscription.Filter.Context,
					subscription.Filter.Group,
					subscription.Options,
					subscription.Created,
				),
		); err != nil {
			return err
		}

		s.postCommitEvent(tx, func() {
			s.callbacks.SubscriptionCreated(subscription.ID)
		})

	}

	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) subscriptionResult(ctx context.Context, row *sql.Rows) (*fftypes.Subscription, error) {
	subscription := fftypes.Subscription{}
	err := row.Scan(
		&subscription.ID,
		&subscription.Namespace,
		&subscription.Name,
		&subscription.Transport,
		&subscription.Filter.Events,
		&subscription.Filter.Topic,
		&subscription.Filter.Context,
		&subscription.Filter.Group,
		&subscription.Options,
		&subscription.Created,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBReadErr, "subscriptions")
	}
	return &subscription, nil
}

func (s *SQLCommon) getSubscriptionEq(ctx context.Context, eq sq.Eq, textName string) (message *fftypes.Subscription, err error) {

	rows, err := s.query(ctx,
		sq.Select(subscriptionColumns...).
			From("subscriptions").
			Where(eq),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Subscription '%s' not found", textName)
		return nil, nil
	}

	subscription, err := s.subscriptionResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SQLCommon) GetSubscriptionByID(ctx context.Context, id *fftypes.UUID) (message *fftypes.Subscription, err error) {
	return s.getSubscriptionEq(ctx, sq.Eq{"id": id}, id.String())
}

func (s *SQLCommon) GetSubscriptionByName(ctx context.Context, ns, name string) (message *fftypes.Subscription, err error) {
	return s.getSubscriptionEq(ctx, sq.Eq{"namespace": ns, "name": name}, fmt.Sprintf("%s:%s", ns, name))
}

func (s *SQLCommon) GetSubscriptions(ctx context.Context, filter database.Filter) (message []*fftypes.Subscription, err error) {

	query, err := s.filterSelect(ctx, "", sq.Select(subscriptionColumns...).From("subscriptions"), filter, subscriptionFilterTypeMap)
	if err != nil {
		return nil, err
	}

	rows, err := s.query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscription := []*fftypes.Subscription{}
	for rows.Next() {
		d, err := s.subscriptionResult(ctx, rows)
		if err != nil {
			return nil, err
		}
		subscription = append(subscription, d)
	}

	return subscription, err

}

func (s *SQLCommon) UpdateSubscription(ctx context.Context, namespace, name string, update database.Update) (err error) {

	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	query, err := s.buildUpdate(sq.Update("subscriptions"), update, subscriptionFilterTypeMap)
	if err != nil {
		return err
	}
	query = query.Where(sq.Eq{
		"namespace": namespace,
		"name":      name,
	})

	err = s.updateTx(ctx, tx, query)
	if err != nil {
		return err
	}

	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) DeleteSubscriptionByID(ctx context.Context, id *fftypes.UUID) (err error) {

	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	err = s.deleteTx(ctx, tx, sq.Delete("subscriptions").Where(sq.Eq{"id": id}))
	if err != nil {
		return err
	}

	s.postCommitEvent(tx, func() {
		s.callbacks.SubscriptionDeleted(id)
	})

	return s.commitTx(ctx, tx, autoCommit)
}
