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
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPinsE2EWithDB(t *testing.T) {
	log.SetLevel("trace")

	s := newQLTestProvider(t)
	defer s.Close()
	ctx := context.Background()

	s.callbacks.On("EventCreated", mock.Anything).Return()

	// Create a new pin entry
	pin := &fftypes.Pin{
		Masked:     true,
		Hash:       fftypes.NewRandB32(),
		Batch:      fftypes.NewUUID(),
		Index:      10,
		Created:    fftypes.Now(),
		Dispatched: false,
	}
	err := s.UpsertPin(ctx, pin)
	assert.NoError(t, err)

	// Query back the pin
	fb := database.PinQueryFactory.NewFilter(ctx)
	filter := fb.And(
		fb.Eq("masked", pin.Masked),
		fb.Eq("hash", pin.Hash),
		fb.Eq("batch", pin.Batch),
		fb.Gt("created", 0),
	)
	pinRes, err := s.GetPins(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pinRes))

	// Set it dispatched
	err = s.SetPinDispatched(ctx, pin.Sequence)
	assert.NoError(t, err)

	// Double insert, checking no error and we keep the dispatched flag
	existingSequence := pin.Sequence
	pin.Sequence = 99999
	err = s.UpsertPin(ctx, pin)
	assert.NoError(t, err)
	pinRes, err = s.GetPins(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pinRes)) // we didn't add twice
	assert.Equal(t, existingSequence, pin.Sequence)
	assert.True(t, pin.Dispatched)

	// Test delete
	err = s.DeletePin(ctx, pin.Sequence)
	assert.NoError(t, err)
	p, err := s.GetPins(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(p))

}

func TestUpsertPinFailBegin(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("pop"))
	err := s.UpsertPin(context.Background(), &fftypes.Pin{})
	assert.Regexp(t, "FF10114", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertPinFailInsert(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	err := s.UpsertPin(context.Background(), &fftypes.Pin{Hash: fftypes.NewRandB32()})
	assert.Regexp(t, "FF10116", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertPinFailExistingSequenceScan(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock.NewRows([]string{"only one"}).AddRow(true))
	mock.ExpectRollback()
	err := s.UpsertPin(context.Background(), &fftypes.Pin{Hash: fftypes.NewRandB32()})
	assert.Regexp(t, "FF10121", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertPinFailCommit(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("pop"))
	err := s.UpsertPin(context.Background(), &fftypes.Pin{Hash: fftypes.NewRandB32()})
	assert.Regexp(t, "FF10119", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPinQueryFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectQuery("SELECT .*").WillReturnError(fmt.Errorf("pop"))
	f := database.PinQueryFactory.NewFilter(context.Background()).Eq("hash", "")
	_, err := s.GetPins(context.Background(), f)
	assert.Regexp(t, "FF10115", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPinBuildQueryFail(t *testing.T) {
	s, _ := newMockProvider().init()
	f := database.PinQueryFactory.NewFilter(context.Background()).Eq("hash", map[bool]bool{true: false})
	_, err := s.GetPins(context.Background(), f)
	assert.Regexp(t, "FF10149.*type", err)
}

func TestGetPinReadMessageFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"pin"}).AddRow("only one"))
	f := database.PinQueryFactory.NewFilter(context.Background()).Eq("hash", "")
	_, err := s.GetPins(context.Background(), f)
	assert.Regexp(t, "FF10121", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSetPinsDispatchedBeginFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("pop"))
	err := s.SetPinDispatched(context.Background(), 12345)
	assert.Regexp(t, "FF10114", err)
}

func TestSetPinsDispatchedUpdateFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	err := s.SetPinDispatched(context.Background(), 12345)
	assert.Regexp(t, "FF10117", err)
}

func TestPinDeleteBeginFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("pop"))
	err := s.DeletePin(context.Background(), 12345)
	assert.Regexp(t, "FF10114", err)
}

func TestPinDeleteFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	err := s.DeletePin(context.Background(), 12345)
	assert.Regexp(t, "FF10118", err)
}
