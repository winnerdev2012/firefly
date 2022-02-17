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

package networkmap

import (
	"context"

	"github.com/hyperledger/firefly/internal/broadcast"
	"github.com/hyperledger/firefly/internal/i18n"
	"github.com/hyperledger/firefly/internal/identity"
	"github.com/hyperledger/firefly/pkg/database"
	"github.com/hyperledger/firefly/pkg/dataexchange"
	"github.com/hyperledger/firefly/pkg/fftypes"
)

type Manager interface {
	RegisterOrganization(ctx context.Context, org *fftypes.Identity, waitConfirm bool) (msg *fftypes.Message, err error)
	RegisterNode(ctx context.Context, waitConfirm bool) (node *fftypes.Identity, msg *fftypes.Message, err error)
	RegisterNodeOrganization(ctx context.Context, waitConfirm bool) (org *fftypes.Identity, msg *fftypes.Message, err error)

	GetOrganizationByID(ctx context.Context, id string) (*fftypes.Identity, error)
	GetOrganizations(ctx context.Context, filter database.AndFilter) ([]*fftypes.Identity, *database.FilterResult, error)
	GetNodeByID(ctx context.Context, id string) (*fftypes.Identity, error)
	GetNodes(ctx context.Context, filter database.AndFilter) ([]*fftypes.Identity, *database.FilterResult, error)
}

type networkMap struct {
	ctx       context.Context
	database  database.Plugin
	broadcast broadcast.Manager
	exchange  dataexchange.Plugin
	identity  identity.Manager
}

func NewNetworkMap(ctx context.Context, di database.Plugin, bm broadcast.Manager, dx dataexchange.Plugin, im identity.Manager) (Manager, error) {
	if di == nil || bm == nil || dx == nil || im == nil {
		return nil, i18n.NewError(ctx, i18n.MsgInitializationNilDepError)
	}

	nm := &networkMap{
		ctx:       ctx,
		database:  di,
		broadcast: bm,
		exchange:  dx,
		identity:  im,
	}
	return nm, nil
}
