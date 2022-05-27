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

package apiserver

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/hyperledger/firefly/mocks/operationmocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSPIPatchOperationByID(t *testing.T) {
	o, r := newTestSPIServer()
	req := httptest.NewRequest("PATCH", "/spi/v1/operations/ns1/abcd12345", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := httptest.NewRecorder()

	mop := &operationmocks.Manager{}
	o.On("Operations").Return(mop)
	mop.On("ResolveOperationByID", mock.Anything, "ns1", "abcd12345", mock.AnythingOfType("*core.OperationUpdateDTO")).Return(nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Result().StatusCode)
}
