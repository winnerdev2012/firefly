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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
	"github.com/stretchr/testify/assert"
)

func TestOperationE2EWithDB(t *testing.T) {

	s := newQLTestProvider(t)
	defer s.Close()
	ctx := context.Background()

	// Create a new operation entry
	operationID := fftypes.NewUUID()
	operation := &fftypes.Operation{
		ID:          operationID,
		Type:        fftypes.OpTypeBlockchainBatchPin,
		Transaction: fftypes.NewUUID(),
		Status:      fftypes.OpStatusPending,
		Created:     fftypes.Now(),
	}
	err := s.UpsertOperation(ctx, operation, true)
	assert.NoError(t, err)

	// Check we get the exact same operation back
	operationRead, err := s.GetOperationByID(ctx, operationID)
	assert.NoError(t, err)
	assert.NotNil(t, operationRead)
	operationJson, _ := json.Marshal(&operation)
	operationReadJson, _ := json.Marshal(&operationRead)
	assert.Equal(t, string(operationJson), string(operationReadJson))

	// Update the operation (this is testing what's possible at the database layer,
	// and does not account for the verification that happens at the higher level)
	operationUpdated := &fftypes.Operation{
		ID:          operationID,
		Type:        fftypes.OpTypeBlockchainBatchPin,
		Transaction: fftypes.NewUUID(),
		Status:      fftypes.OpStatusFailed,
		Recipient:   "sally",
		Plugin:      "ethereum",
		BackendID:   fftypes.NewRandB32().String(),
		Error:       "pop",
		Info:        fftypes.JSONObject{"some": "info"},
		Created:     fftypes.Now(),
		Updated:     fftypes.Now(),
	}
	err = s.UpsertOperation(context.Background(), operationUpdated, true)
	assert.NoError(t, err)

	// Check we get the exact same message back - note the removal of one of the operation elements
	operationRead, err = s.GetOperationByID(ctx, operationID)
	assert.NoError(t, err)
	operationJson, _ = json.Marshal(&operationUpdated)
	operationReadJson, _ = json.Marshal(&operationRead)
	assert.Equal(t, string(operationJson), string(operationReadJson))

	// Query back the operation
	fb := database.OperationQueryFactory.NewFilter(ctx)
	filter := fb.And(
		fb.Eq("id", operationUpdated.ID.String()),
		fb.Eq("tx", operationUpdated.Transaction),
		fb.Eq("type", operationUpdated.Type),
		fb.Eq("recipient", operationUpdated.Recipient),
		fb.Eq("status", operationUpdated.Status),
		fb.Eq("error", operationUpdated.Error),
		fb.Eq("plugin", operationUpdated.Plugin),
		fb.Eq("backendid", operationUpdated.BackendID),
		fb.Gt("created", 0),
		fb.Gt("updated", 0),
	)

	operations, err := s.GetOperations(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(operations))
	operationReadJson, _ = json.Marshal(operations[0])
	assert.Equal(t, string(operationJson), string(operationReadJson))

	// Negative test on filter
	filter = fb.And(
		fb.Eq("id", operationUpdated.ID.String()),
		fb.Eq("updated", "0"),
	)
	operations, err = s.GetOperations(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(operations))

	// Update
	updateTime := fftypes.Now()
	up := database.OperationQueryFactory.NewUpdate(ctx).
		Set("status", fftypes.OpStatusSucceeded).
		Set("updated", updateTime).
		Set("error", "")
	idFilter := database.OperationQueryFactory.NewFilter(ctx).
		Eq("id", operationUpdated.ID)
	err = s.UpdateOperations(ctx, idFilter, up)
	assert.NoError(t, err)

	// Test find updated value
	filter = fb.And(
		fb.Eq("id", operationUpdated.ID.String()),
		fb.Eq("status", fftypes.OpStatusSucceeded),
		fb.Eq("updated", updateTime),
		fb.Eq("error", ""),
	)
	operations, err = s.GetOperations(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(operations))
}

func TestUpsertOperationFailBegin(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("pop"))
	err := s.UpsertOperation(context.Background(), &fftypes.Operation{}, true)
	assert.Regexp(t, "FF10114", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertOperationFailSelect(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	operationID := fftypes.NewUUID()
	err := s.UpsertOperation(context.Background(), &fftypes.Operation{ID: operationID}, true)
	assert.Regexp(t, "FF10115", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertOperationFailInsert(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{}))
	mock.ExpectExec("INSERT .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	operationID := fftypes.NewUUID()
	err := s.UpsertOperation(context.Background(), &fftypes.Operation{ID: operationID}, true)
	assert.Regexp(t, "FF10116", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertOperationFailUpdate(t *testing.T) {
	s, mock := newMockProvider().init()
	operationID := fftypes.NewUUID()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(operationID.String()))
	mock.ExpectExec("UPDATE .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	err := s.UpsertOperation(context.Background(), &fftypes.Operation{ID: operationID}, true)
	assert.Regexp(t, "FF10117", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpsertOperationFailCommit(t *testing.T) {
	s, mock := newMockProvider().init()
	operationID := fftypes.NewUUID()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec("INSERT .*").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("pop"))
	err := s.UpsertOperation(context.Background(), &fftypes.Operation{ID: operationID}, true)
	assert.Regexp(t, "FF10119", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOperationByIDSelectFail(t *testing.T) {
	s, mock := newMockProvider().init()
	operationID := fftypes.NewUUID()
	mock.ExpectQuery("SELECT .*").WillReturnError(fmt.Errorf("pop"))
	_, err := s.GetOperationByID(context.Background(), operationID)
	assert.Regexp(t, "FF10115", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOperationByIDNotFound(t *testing.T) {
	s, mock := newMockProvider().init()
	operationID := fftypes.NewUUID()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"id"}))
	msg, err := s.GetOperationByID(context.Background(), operationID)
	assert.NoError(t, err)
	assert.Nil(t, msg)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOperationByIDScanFail(t *testing.T) {
	s, mock := newMockProvider().init()
	operationID := fftypes.NewUUID()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("only one"))
	_, err := s.GetOperationByID(context.Background(), operationID)
	assert.Regexp(t, "FF10121", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOperationsQueryFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectQuery("SELECT .*").WillReturnError(fmt.Errorf("pop"))
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", "")
	_, err := s.GetOperations(context.Background(), f)
	assert.Regexp(t, "FF10115", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOperationsBuildQueryFail(t *testing.T) {
	s, _ := newMockProvider().init()
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", map[bool]bool{true: false})
	_, err := s.GetOperations(context.Background(), f)
	assert.Regexp(t, "FF10149.*id", err)
}

func TestGettOperationsReadMessageFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectQuery("SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("only one"))
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", "")
	_, err := s.GetOperations(context.Background(), f)
	assert.Regexp(t, "FF10121", err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestOperationUpdateBeginFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin().WillReturnError(fmt.Errorf("pop"))
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", fftypes.NewUUID())
	u := database.OperationQueryFactory.NewUpdate(context.Background()).Set("id", fftypes.NewUUID())
	err := s.UpdateOperations(context.Background(), f, u)
	assert.Regexp(t, "FF10114", err)
}

func TestOperationUpdateBuildQueryFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", fftypes.NewUUID())
	u := database.OperationQueryFactory.NewUpdate(context.Background()).Set("id", map[bool]bool{true: false})
	err := s.UpdateOperations(context.Background(), f, u)
	assert.Regexp(t, "FF10149.*id", err)
}

func TestOperationUpdateBuildFilterFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", map[bool]bool{true: false})
	u := database.OperationQueryFactory.NewUpdate(context.Background()).Set("id", fftypes.NewUUID())
	err := s.UpdateOperations(context.Background(), f, u)
	assert.Regexp(t, "FF10149.*id", err)
}

func TestOperationUpdateFail(t *testing.T) {
	s, mock := newMockProvider().init()
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE .*").WillReturnError(fmt.Errorf("pop"))
	mock.ExpectRollback()
	f := database.OperationQueryFactory.NewFilter(context.Background()).Eq("id", fftypes.NewUUID())
	u := database.OperationQueryFactory.NewUpdate(context.Background()).Set("id", fftypes.NewUUID())
	err := s.UpdateOperations(context.Background(), f, u)
	assert.Regexp(t, "FF10117", err)
}
