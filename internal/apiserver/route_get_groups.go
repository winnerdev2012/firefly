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

	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly/internal/coremsgs"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/database"
)

var getGroups = &ffapi.Route{
	Name:            "getGroups",
	Path:            "groups",
	Method:          http.MethodGet,
	PathParams:      nil,
	QueryParams:     nil,
	Description:     coremsgs.APIEndpointsGetGroups,
	JSONInputValue:  nil,
	JSONOutputValue: func() interface{} { return []*core.Group{} },
	JSONOutputCodes: []int{http.StatusOK},
	Extensions: &coreExtensions{
		FilterFactory: database.GroupQueryFactory,
		CoreJSONHandler: func(r *ffapi.APIRequest, cr *coreRequest) (output interface{}, err error) {
			pm := cr.or.PrivateMessaging()
			if pm == nil {
				return nil, i18n.NewError(cr.ctx, coremsgs.MsgActionOnlyValidMultiparty)
			}
			return filterResult(pm.GetGroups(cr.ctx, cr.filter))
		},
	},
}
