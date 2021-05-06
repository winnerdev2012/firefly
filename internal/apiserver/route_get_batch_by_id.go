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

package apiserver

import (
	"net/http"

	"github.com/kaleido-io/firefly/internal/engine"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
)

var getBatchById = &Route{
	Name:   "getBatchById",
	Path:   "ns/{ns}/batches/{id}",
	Method: http.MethodGet,
	PathParams: []PathParam{
		{Name: "ns", Description: i18n.MsgTBD},
		{Name: "id", Description: i18n.MsgTBD},
	},
	QueryParams:     nil,
	Description:     i18n.MsgTBD,
	JSONInputValue:  func() interface{} { return nil },
	JSONOutputValue: func() interface{} { return &fftypes.Batch{} },
	JSONHandler: func(e engine.Engine, req *http.Request, pp map[string]string, qp map[string]string, input interface{}) (output interface{}, status int, err error) {
		output, err = e.GetBatchById(req.Context(), pp["ns"], pp["id"])
		return output, 200, err
	},
}
