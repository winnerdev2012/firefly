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
	"net/http"

	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/oapispec"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

var getDataDefByID = &oapispec.Route{
	Name:   "getDataDefByID",
	Path:   "namespaces/{ns}/datatypes/{dtypeid}",
	Method: http.MethodGet,
	PathParams: []*oapispec.PathParam{
		{Name: "ns", ExampleFromConf: config.NamespacesDefault, Description: i18n.MsgTBD},
		{Name: "dtypeid", Description: i18n.MsgTBD},
	},
	QueryParams:     nil,
	FilterFactory:   nil,
	Description:     i18n.MsgTBD,
	JSONInputValue:  nil,
	JSONOutputValue: func() interface{} { return &fftypes.Datatype{} },
	JSONOutputCode:  http.StatusOK,
	JSONHandler: func(r oapispec.APIRequest) (output interface{}, err error) {
		output, err = r.Or.GetDatatypeByID(r.Ctx, r.PP["ns"], r.PP["dtypeid"])
		return output, err
	},
}
