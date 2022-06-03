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
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
)

var getMsgByID = &ffapi.Route{
	Name:   "getMsgByID",
	Path:   "messages/{msgid}",
	Method: http.MethodGet,
	PathParams: []*ffapi.PathParam{
		{Name: "msgid", Description: coremsgs.APIParamsMessageID},
	},
	QueryParams: []*ffapi.QueryParam{
		{Name: "fetchdata", IsBool: true, Description: coremsgs.APIFetchDataDesc},
	},
	Description:     coremsgs.APIEndpointsGetMsgByID,
	JSONInputValue:  nil,
	JSONOutputValue: func() interface{} { return &core.MessageInOut{} }, // can include full values
	JSONOutputCodes: []int{http.StatusOK},
	Extensions: &coreExtensions{
		CoreJSONHandler: func(r *ffapi.APIRequest, cr *coreRequest) (output interface{}, err error) {
			if strings.EqualFold(r.QP["data"], "true") || strings.EqualFold(r.QP["fetchdata"], "true") {
				return cr.or.GetMessageByIDWithData(cr.ctx, extractNamespace(r.PP), r.PP["msgid"])
			}
			return cr.or.GetMessageByID(cr.ctx, extractNamespace(r.PP), r.PP["msgid"])
		},
	},
}
