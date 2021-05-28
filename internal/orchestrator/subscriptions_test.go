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
	"context"
	"fmt"
	"testing"

	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func uuidMatches(id1 *fftypes.UUID) interface{} {
	return mock.MatchedBy(func(id2 *fftypes.UUID) bool { return *id1 == *id2 })
}

func TestCreateSubscriptionBadNamespace(t *testing.T) {
	or := newTestOrchestrator()
	_, err := or.CreateSubscription(or.ctx, "!wrong", &fftypes.Subscription{
		SubscriptionRef: fftypes.SubscriptionRef{
			Name: "sub1",
		},
	})
	assert.Regexp(t, "FF10131", err)
}

func TestCreateSubscriptionNamespace(t *testing.T) {
	or := newTestOrchestrator()
	or.mdi.On("GetNamespace", mock.Anything, "ns1").Return(nil, nil)
	_, err := or.CreateSubscription(or.ctx, "ns1", &fftypes.Subscription{
		SubscriptionRef: fftypes.SubscriptionRef{
			Name: "sub1",
		},
	})
	assert.Regexp(t, "FF10187", err)
}

func TestCreateSubscriptionBadName(t *testing.T) {
	or := newTestOrchestrator()
	or.mdi.On("GetNamespace", mock.Anything, "ns1").Return(&fftypes.Namespace{}, nil)
	_, err := or.CreateSubscription(or.ctx, "ns1", &fftypes.Subscription{
		SubscriptionRef: fftypes.SubscriptionRef{
			Name: "!sub1",
		},
	})
	assert.Regexp(t, "FF10131", err)
}

func TestCreateSubscriptionOk(t *testing.T) {
	or := newTestOrchestrator()
	sub := &fftypes.Subscription{
		SubscriptionRef: fftypes.SubscriptionRef{
			Name: "sub1",
		},
	}
	or.mdi.On("GetNamespace", mock.Anything, "ns1").Return(&fftypes.Namespace{}, nil)
	or.mei.On("CreateDurableSubscription", mock.Anything, mock.Anything).Return(nil)
	s1, err := or.CreateSubscription(or.ctx, "ns1", sub)
	assert.NoError(t, err)
	assert.Equal(t, s1, sub)
	assert.Equal(t, "ns1", sub.Namespace)
}

func TestDeleteSubscriptionBadUUID(t *testing.T) {
	or := newTestOrchestrator()
	or.mdi.On("GetSubscriptionByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("pop"))
	err := or.DeleteSubscription(or.ctx, "ns2", "! a UUID")
	assert.Regexp(t, "FF10142", err)
}

func TestDeleteSubscriptionLookupError(t *testing.T) {
	or := newTestOrchestrator()
	or.mdi.On("GetSubscriptionByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("pop"))
	err := or.DeleteSubscription(or.ctx, "ns2", fftypes.NewUUID().String())
	assert.EqualError(t, err, "pop")
}

func TestDeleteSubscriptionNSMismatch(t *testing.T) {
	or := newTestOrchestrator()
	sub := &fftypes.Subscription{
		SubscriptionRef: fftypes.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Name:      "sub1",
			Namespace: "ns1",
		},
	}
	or.mdi.On("GetSubscriptionByID", mock.Anything, sub.ID).Return(sub, nil)
	err := or.DeleteSubscription(or.ctx, "ns2", sub.ID.String())
	assert.Regexp(t, "FF10109", err)
}

func TestDeleteSubscription(t *testing.T) {
	or := newTestOrchestrator()
	sub := &fftypes.Subscription{
		SubscriptionRef: fftypes.SubscriptionRef{
			ID:        fftypes.NewUUID(),
			Name:      "sub1",
			Namespace: "ns1",
		},
	}
	or.mdi.On("GetSubscriptionByID", mock.Anything, uuidMatches(sub.ID)).Return(sub, nil)
	or.mei.On("DeleteDurableSubscription", mock.Anything, sub).Return(nil)
	err := or.DeleteSubscription(or.ctx, "ns1", sub.ID.String())
	assert.NoError(t, err)
}

func TestGetSubscriptions(t *testing.T) {
	or := newTestOrchestrator()
	u := fftypes.NewUUID()
	or.mdi.On("GetSubscriptions", mock.Anything, mock.Anything).Return([]*fftypes.Subscription{}, nil)
	fb := database.SubscriptionQueryFactory.NewFilter(context.Background())
	f := fb.And(fb.Eq("id", u))
	_, err := or.GetSubscriptions(context.Background(), "ns1", f)
	assert.NoError(t, err)
}

func TestGetSGetSubscriptionsByID(t *testing.T) {
	or := newTestOrchestrator()
	u := fftypes.NewUUID()
	or.mdi.On("GetSubscriptionByID", mock.Anything, u).Return(nil, nil)
	_, err := or.GetSubscriptionByID(context.Background(), "ns1", u.String())
	assert.NoError(t, err)
}

func TestGetSubscriptionDefsByIDBadID(t *testing.T) {
	or := newTestOrchestrator()
	_, err := or.GetSubscriptionByID(context.Background(), "", "")
	assert.Regexp(t, "FF10142", err)
}
