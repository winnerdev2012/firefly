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

package apiserver

import (
	"net/http"
	"strings"

	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/internal/orchestrator"
)

var getContractInterfaceNameVersion = &ffapi.Route{
	Name:   "getContractInterfaceByNameAndVersion",
	Path:   "contracts/interfaces/{name}/{version}",
	Method: http.MethodGet,
	PathParams: []*ffapi.PathParam{
		{Name: "name", Description: coremsgs.APIParamsContractInterfaceName},
		{Name: "version", Description: coremsgs.APIParamsContractInterfaceVersion},
	},
	QueryParams: []*ffapi.QueryParam{
		{Name: "fetchchildren", Example: "true", Description: coremsgs.APIParamsContractInterfaceFetchChildren, IsBool: true},
	},
	Description:     coremsgs.APIEndpointsGetContractInterfaceNameVersion,
	JSONInputValue:  nil,
	JSONOutputValue: func() interface{} { return &fftypes.FFI{} },
	JSONOutputCodes: []int{http.StatusOK},
	Extensions: &coreExtensions{
		EnabledIf: func(or orchestrator.Orchestrator) bool {
			return or.Contracts() != nil
		},
		CoreJSONHandler: func(r *ffapi.APIRequest, cr *coreRequest) (output interface{}, err error) {
			if strings.EqualFold(r.QP["fetchchildren"], "true") {
				return cr.or.Contracts().GetFFIWithChildren(cr.ctx, r.PP["name"], r.PP["version"])
			}
			return cr.or.Contracts().GetFFI(cr.ctx, r.PP["name"], r.PP["version"])
		},
	},
}
