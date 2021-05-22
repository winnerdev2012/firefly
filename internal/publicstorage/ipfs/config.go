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

package ipfs

import (
	"github.com/kaleido-io/firefly/internal/config"
	"github.com/kaleido-io/firefly/internal/restclient"
)

const (
	// IPFSConfAPISubconf is the http configuration to connect to the API endpoint of IPFS
	IPFSConfAPISubconf = "api"
	// IPFSConfGatewaySubconf is the http configuration to connect to the Gateway endpoint of IPFS
	IPFSConfGatewaySubconf = "gateway"
)

func (i *IPFS) InitPrefix(prefix config.Prefix) {
	restclient.InitPrefix(prefix.SubPrefix(IPFSConfAPISubconf))
	restclient.InitPrefix(prefix.SubPrefix(IPFSConfGatewaySubconf))
}
